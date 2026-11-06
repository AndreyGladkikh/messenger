package outbox_relay

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	cnvEvnt "messenger/messenger/internal/messaging/app/event"
	"messenger/messenger/internal/platform/event"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/rabbitmq/producers"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/platform/transactions"
	appEvent "messenger/messenger/internal/shared/app/event"
	"strings"
	"time"

	"uuid"
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
	Logger       *logger.Logger
	EventService *event.EventService
	TXManager    *transactions.Manager
	Producers    *producers.Producers
}

func New(
	Logger *logger.Logger,
	EventService *event.EventService,
	TXManager *transactions.Manager,
	Producers *producers.Producers,
) *OutboxRelay {
	return &OutboxRelay{
		Logger:       Logger,
		EventService: EventService,
		TXManager:    TXManager,
		Producers:    Producers,
	}
}

func (r *OutboxRelay) Run(ctx context.Context) error {
	r.Logger.InfoContext(ctx, "outbox relay started")

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		r.Logger.InfoContext(ctx, "outbox relay stopped")
		<-time.After(30 * time.Second)
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := r.run(runCtx); err != nil {
				r.Logger.ErrorContext(ctx, "outbox relay process failed", "error", err)
				return fmt.Errorf("outbox relay process failed: %w", err)
			}
		}
	}
}

func (r *OutboxRelay) run(ctx context.Context) error {
	var err error

	outboxEvent, err := r.EventService.GetNextUnprocessedEvent(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to retreive event to process: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		time.Sleep(1 * time.Second)
		return nil
	}

	// todo refactor, rm
	// eventEnvelope, err := event.TranslateOutboxMsgToEventEnvelope(outboxEvent)
	// if err != nil {
	// 	return fmt.Errorf("%w: failed to translate stored event to dispatched event: %w", ErrFatal, err)
	// }

	err = nil
	producer := r.Producer(outboxEvent)
	if producer != nil {
		envlp := event.Envelope{
			ID:         outboxEvent.ID,
			Type:       outboxEvent.Type,
			OccurredAt: outboxEvent.OccurredAt,
			Event:      outboxEvent.Payload,
		}

		err = producer.Produce(ctx, envlp.Type, envlp)
	} else {
		// todo dl
		err = errors.New("no producer for outbox")
	}

	var (
		status      event.Status
		nextRetryAt *time.Time
	)

	if err != nil {
		if errors.Is(err, ErrFatal) || outboxEvent.Attempts+1 >= retries {
			status, nextRetryAt = event.StatusDead, nil
		} else {
			retr := r.calcNextRetry()
			status, nextRetryAt = event.StatusRetry, &retr
		}

		r.EventService.UpdateEvent(ctx, outboxEvent, status, err, nextRetryAt)
		return nil
	}

	r.EventService.UpdateEvent(ctx, outboxEvent, event.StatusSucceeded, err, nil)
	return nil
}

// func (r *OutboxRelay) executeEventHandlers(ctx context.Context, outboxEvent sqlc.Outbox) error {
// 	eventEnvelope, err := event.TranslateOutboxMsgToEventEnvelope(outboxEvent)
// 	if err != nil {
// 		return fmt.Errorf("%w: failed to translate stored event to dispatched event: %w", ErrFatal, err)
// 	}

// 	eventHandlers := r.EventHandlerRegistry.HandlersForEvent(eventEnvelope.Event)

// 	var errs []error
// 	var wg sync.WaitGroup
// 	for _, handler := range eventHandlers {
// 		wg.Go(func() {
// 			if err := r.runHandler(ctx, outboxEvent.EventID, eventEnvelope, handler); err != nil {
// 				errs = append(errs, fmt.Errorf("failed to execute handler %q: %w", handler.Name(), err))
// 			}
// 		})
// 	}
// 	wg.Wait()
// 	if len(errs) > 0 {
// 		return &EventHandlingError{errs}
// 	}
// 	return nil
// }

// func (r *OutboxRelay) executeEventHandlersV2(ctx context.Context, outboxEvent sqlc.Outbox) error {
// 	eventEnvelope, err := event.TranslateOutboxMsgToEventEnvelope(outboxEvent)
// 	if err != nil {
// 		return fmt.Errorf("%w: failed to translate stored event to dispatched event: %w", ErrFatal, err)
// 	}

// 	eventHandlers := handlersRegistry[eventEnvelope.Event.Name()]

// 	var errs []error
// 	var wg sync.WaitGroup
// 	for _, handler := range eventHandlers {
// 		wg.Go(func() {
// 			if err := r.runHandler(ctx, outboxEvent.EventID, eventEnvelope, handler); err != nil {
// 				errs = append(errs, fmt.Errorf("failed to execute handler %q: %w", handler.Name(), err))
// 			}
// 		})
// 	}
// 	wg.Wait()
// 	if len(errs) > 0 {
// 		return &EventHandlingError{errs}
// 	}
// 	return nil
// }

func (r *OutboxRelay) calcNextRetry() time.Time {
	return time.Now().Add(5 * time.Second)
}

func (r *OutboxRelay) runHandler(ctx context.Context, eventID uuid.UUID, envelope appEvent.Envelope[appEvent.Event], handler event.Handler) error {
	return r.TXManager.WithTransaction(ctx, func(ctx context.Context) error {
		err := r.EventService.RegisterEventHandlerExecution(ctx, eventID, handler.Name())
		if err != nil && !errors.Is(err, event.ErrHandlerAlreadyExecuted) {
			return err
		}
		if errors.Is(err, event.ErrHandlerAlreadyExecuted) {
			return nil
		}
		return handler.Handle(ctx, envelope)
	})
}

// func (r *OutboxRelay) runHandlerV2(ctx context.Context, eventID uuid.UUID, envelope apevent.Envelope[domain.Event], handler event.Handler) error {
// 	return r.txManager.WithTransaction(ctx, func(ctx context.Context) error {
// 		err := r.eventService.RegisterEventHandlerExecution(ctx, eventID, handler.Name())
// 		if err != nil && !errors.Is(err, event.ErrHandlerAlreadyExecuted) {
// 			return err
// 		}
// 		if errors.Is(err, event.ErrHandlerAlreadyExecuted) {
// 			return nil
// 		}
// 		return handler.Handle(ctx, envelope)
// 	})
// }

// var handlersRegistry = map[string][]Handler[any] {
// 	entry.EntryRegisteredEventName: {
// 		DirectoryEntryRegisteredSearchCreateIndexHandler,
// 	},
// }

// type Handler[E any] interface {
// 	Name() string
// 	Handle(ctx context.Context, eventEnvlp apevent.Envelope[E]) error
// }

// type DirectoryEntryRegisteredSearchCreateIndexHandler struct {
// 	CreateGlobalSearchIndexHandler *srchProjection.CreateGlobalSearchIndexHandler
// }

// func NewDirectoryEntryRegisteredSearchCreateIndexHandler(
// 	CreateGlobalSearchIndexHandler *srchProjection.CreateGlobalSearchIndexHandler,
// ) *DirectoryEntryRegisteredSearchCreateIndexHandler {
// 	return &DirectoryEntryRegisteredSearchCreateIndexHandler{
// 		CreateGlobalSearchIndexHandler: CreateGlobalSearchIndexHandler,
// 	}
// }

// func (*DirectoryEntryRegisteredSearchCreateIndexHandler) Name() string {
// 	return "bc:search.event:directory_entry_registered.handler:create_index"
// }

// func (h *DirectoryEntryRegisteredSearchCreateIndexHandler) Handle(ctx context.Context, eventEnvlp apevent.Envelope[entry.EntryRegistered]) error {
// 	err := h.CreateGlobalSearchIndexHandler.Handle(ctx, srchProjection.CreateGlobalSearchIndex{
// 		ID: eventEnvlp.Event.TargetID,
// 		Type: eventEnvlp.Event.TargetType.String(),
// 		DisplayName: eventEnvlp.Event.DisplayName,
// 		Handle: eventEnvlp.Event.Handle,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

type Producer interface {
	Produce(ctx context.Context, typ string, payload any) error
}

func (o *OutboxRelay) Producer(outbox sqlc.Outbox) Producer {
	switch outbox.Type {
	case cnvEvnt.DirectConversationStartedEventName:
		switch outbox.Kind {
		case sqlc.OutboxKindIntegration:
			return o.Producers.ConversationIntegrationEventsProducer
		}
	}

	return nil
}
