package outbox_relay

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/messenger/internal/messaging/domain"
	"messenger/messenger/internal/messaging/infrastructure/event"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/db/transaction"
	"messenger/messenger/internal/platform/logger"
	"strings"
	"sync"
	"time"
)

var ErrFatal = errors.New("fatal error")

const retries = 5

type EventHandlingError struct {
	handlerErrors []error
}

func (e *EventHandlingError) Error() string {
	return strings.Join(e.AsStrings(), "; ")
}

func (e *EventHandlingError) AsStrings() []string {
	s := make([]string, 0, len(e.handlerErrors))
	for _, e := range e.handlerErrors {
		s = append(s, e.Error())
	}
	return s
}

type OutboxRelay struct {
	eventService         *event.EventService
	txManager            *transaction.Manager
	logger               *logger.Logger
	qs                   *sqlc.Queries
	eventHandlerRegistry *event.HandlerRegistry
}

func NewOutboxRelay(
	eventService *event.EventService,
	txManager *transaction.Manager,
	logger *logger.Logger,
	qs *sqlc.Queries,
	eventHandlerRegistry *event.HandlerRegistry,
) *OutboxRelay {
	return &OutboxRelay{
		eventService:         eventService,
		txManager:            txManager,
		logger:               logger,
		qs:                   qs,
		eventHandlerRegistry: eventHandlerRegistry,
	}
}

func (r *OutboxRelay) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := r.processEvent(); err != nil {
				r.logger.Error("outbox relay error", "error", err)
				return err
			}
		}
	}
}

func (r *OutboxRelay) processEvent() error {
	ctx := context.Background()

	outboxEvent, err := r.eventService.GetNextUnprocessedEvent(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to retreive event to process: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		time.Sleep(1 * time.Second)
		return nil
	}

	var (
		status      event.Status
		nextRetryAt time.Time
		errs        []error
	)
	err = r.executeEventHandlers(ctx, outboxEvent)
	if err != nil {
		if errors.Is(err, ErrFatal) || outboxEvent.Attempts+1 >= retries {
			status, nextRetryAt = event.StatusDead, time.Time{}
		} else {
			status, nextRetryAt = event.StatusRetry, r.calcNextRetry()
		}

		if ehErr, ok := errors.AsType[*EventHandlingError](err); ok {
			errs = ehErr.handlerErrors
		} else {
			errs = append(errs, err)
		}
		r.eventService.UpdateEvent(ctx, outboxEvent, status, errs, nextRetryAt)
		return nil
	}

	r.eventService.UpdateEvent(ctx, outboxEvent, event.StatusSucceeded, errs, time.Time{})
	return nil
}

func (r *OutboxRelay) executeEventHandlers(ctx context.Context, outboxEvent sqlc.Outbox) error {
	domainEvent, err := translateStoredEventToDomainEvent(outboxEvent)
	if err != nil {
		return fmt.Errorf("%w: failed to translate stored event to dispatched event: %w", ErrFatal, err)
	}

	eventHandlers := r.eventHandlerRegistry.HandlersForEvent(domainEvent)

	var errs []error
	var wg sync.WaitGroup
	for _, handler := range eventHandlers {
		wg.Go(func() {
			if err := r.runHandler(ctx, outboxEvent.EventID, domainEvent, handler); err != nil {
				errs = append(errs, fmt.Errorf("failed to execute handler %q: %w", handler.Name(), err))
			}
		})
	}
	wg.Wait()
	if len(errs) > 0 {
		return &EventHandlingError{errs}
	}
	return nil
}

func (r *OutboxRelay) calcNextRetry() time.Time {
	return time.Now().Add(5 * time.Second)
}

func (r *OutboxRelay) runHandler(ctx context.Context, eventID string, domainEvent domain.Event, handler event.Handler) error {
	return r.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		err := r.eventService.RegisterEventHandlerExecution(ctx, eventID, handler.Name())
		if err != nil && !errors.Is(err, event.ErrHandlerAlreadyExecuted) {
			return err
		}
		if errors.Is(err, event.ErrHandlerAlreadyExecuted) {
			return nil
		}
		return handler.Handle(ctx, domainEvent)
	})
}
