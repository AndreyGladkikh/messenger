package event

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/infrastructure/db"
	"messenger/messenger/internal/infrastructure/logger"
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

type Processor struct {
	eventStorage *EventStorage
	txManager    *transaction.Manager
	logger       *logger.Logger
	qs           *db.Queries
	eventHandlerRegistry *HandlerRegistry
}

func NewProcessor(
	eventStorage *EventStorage,
	txManager *transaction.Manager,
	logger *logger.Logger,
	qs *db.Queries,
	eventHandlerRegistry *HandlerRegistry,
) *Processor {
	return &Processor{
		eventStorage: eventStorage,
		txManager:    txManager,
		logger:       logger,
		qs:           qs,
		eventHandlerRegistry: eventHandlerRegistry,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := p.processEvent(); err != nil {
				return err
			}
		}
	}
}

func (p *Processor) processEvent() error {
	ctx := context.Background()

	outboxEvent, err := p.eventStorage.GetNextUnprocessedEvent(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("retreive event to process: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		time.Sleep(1 * time.Second)
		return nil
	}

	var (
		status status
		nextRetryAt time.Time
		errs []error
	)
	err = p.executeEventHandlers(ctx, outboxEvent)
	if err != nil {
		if errors.Is(err, ErrFatal) || outboxEvent.Attempts + 1 >= retries {
			status, nextRetryAt = statusDead, time.Time{}
		} else {
			status, nextRetryAt = statusRetry, p.calcNextRetry()
		}
		
		if ehErr, ok := errors.AsType[*EventHandlingError](err); ok {
			errs = ehErr.handlerErrors
		} else {
			errs = append(errs, err)
		}
		p.eventStorage.UpdateEvent(ctx, outboxEvent, status, errs, nextRetryAt)
		return nil
	}

	p.eventStorage.UpdateEvent(ctx, outboxEvent, statusSucceeded, errs, time.Time{})
	return nil
}

func (p *Processor) executeEventHandlers(ctx context.Context, outboxEvent db.Outbox) error {
	domainEvent, err := translateStoredEventToDomainEvent(outboxEvent)
	if err != nil {
		return fmt.Errorf("%w: failed to translate stored event to dispatched event: %w", ErrFatal, err)
	}

	eventHandlers := p.eventHandlerRegistry.HandlersForEvent(domainEvent)

	var errs []error
	var wg sync.WaitGroup
	for _, handler := range eventHandlers {
		wg.Go(func() {
			if err := p.runHandler(ctx, outboxEvent.EventID, domainEvent, handler); err != nil {
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

func (p *Processor) calcNextRetry() time.Time {
	return time.Now().Add(5*time.Second)
}

func (p *Processor) runHandler(ctx context.Context, eventID string, domainEvent domain.Event, handler Handler) error {
	return p.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		err := p.eventStorage.RegisterEventHandlerExecution(ctx, eventID, handler.Name())
		if err != nil && !errors.Is(err, ErrHandlerAlreadyExecuted) {
			return err
		}
		if errors.Is(err, ErrHandlerAlreadyExecuted) {
			return nil
		}
		return handler.Handle(ctx, domainEvent)
	})
}