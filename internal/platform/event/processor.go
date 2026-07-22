package event

import (
	"context"
	"database/sql"
	"fmt"
	"messenger/messenger/internal/adapters/out/logger"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/application/event"
	"messenger/messenger/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ProcessorConfig struct {
}

type Processor struct {
	logger    *logger.Logger
	qs        *queries.Queries
	eventBus  *Bus[domain.Event]
	cfg       *ProcessorConfig
	inShutdown bool
}

func NewProcessor(
	logger *logger.Logger,
	qs *queries.Queries,
	eventBus *Bus[domain.Event],
) *Processor {
	return &Processor{
		logger:    logger,
		qs:        qs,
		eventBus:  eventBus,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	for {
		if p.inShutdown {
			return nil
		}
		events, err := p.qs.ListUnprocessedEvents(ctx)
		if err != nil {
			return fmt.Errorf("event processor: failed to retreive events: %w", err)
		}
		// todo sleep 1sec if empty

		for _, e := range events {
			go func() {
				err := p.processEvent(ctx, e)
				if err != nil {
					p.logger.Error(
						"event processor: failed to process event",
						"eventID", e.ID,
						"error", err.Error(),
					)
				}
			}()
			err := p.qs.ProcessEvent(ctx, queries.ProcessEventParams{
				ID: e.ID,
				ProcessedAt: mapping.ToDBTimestamp(time.Now()),
			})
			if err != nil {
				return fmt.Errorf("event processor: failed to mark event as processed: %w", err)
			}
		}
	}
}

func (p *Processor) Shutdown(ctx context.Context) error {
	p.inShutdown = true
	return nil
}

func (p *Processor) processEvent(ctx context.Context, storedEvent queries.Event) error {
	event, err := translateStoredEventToDispatchedEvent(storedEvent)
	if err != nil {
		return fmt.Errorf("event processor: failed to translate stored event to dispatched event")
	}
	handlers := p.eventBus.HandlersForEvent(event)
	handlerExecutions, err := p.qs.GetEventHandlerExecutionsByEventId(ctx, storedEvent.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("event processor: failed to retrieve event handler executions: %w", err)
	}
	handlerExecutionsMap := make(map[string]queries.EventHandlerExecution)
	for _, execution := range handlerExecutions {
		handlerExecutionsMap[execution.HandlerType] = execution
	}

	for _, h := range handlers {
		if !shouldExecuteHandler(h, handlerExecutionsMap) {
			continue
		}

		err := p.executeHandler(ctx, storedEvent, event, h, handlerExecutionsMap)
		if err != nil {
			return fmt.Errorf("event processor: failed to execute handler %q: %w", h.Name(), err)
		}
	}

	return nil
}

func (p *Processor) executeHandler(ctx context.Context, storedEvent queries.Event, event domain.Event, handler event.Handler[domain.Event], handlerExecutions map[string]queries.EventHandlerExecution) error {
	execution, err := p.getExecution(ctx, storedEvent, handler, handlerExecutions)
	if err != nil {
		return err
	}

	err = p.eventBus.DispatchForHandler(ctx, event, handler)
	if err != nil {
		err = fmt.Errorf("failed to execute event handler: %w", err)
	}

	updateParams := queries.UpdateEventHandlerExecutionParams{
		ID:       execution.ID,
		Attempts: execution.Attempts + 1,
	}
	if err != nil {
		updateParams.Error = mapping.PgText(err.Error())
		updateParams.NextRetryAt = mapping.ToDBTimestamp(time.Now().Add(5 * time.Minute))
	} else {
		updateParams.Error = mapping.PgText("")
		updateParams.NextRetryAt = mapping.ToDBTimestamp(time.Time{})
	}
	if e := p.qs.UpdateEventHandlerExecution(ctx, updateParams); e != nil {
		err = fmt.Errorf("failed to update event handler execution: %w; %w", e, err)
	}
	return err
}

func shouldExecuteHandler[T domain.Event](handler event.Handler[T], handlerExecutions map[string]queries.EventHandlerExecution) bool {
	var execution queries.EventHandlerExecution

	execution, ok := handlerExecutions[handler.Name()]
	if !ok {
		return true
	}

	if execution.NextRetryAt.Valid && time.Now().Before(execution.NextRetryAt.Time) {
		return false
	}
	if execution.Attempts > 0 && !execution.Error.Valid {
		return false
	}
	if execution.Attempts >= 3 {
		return false
	}

	return true
}

func (p *Processor) getExecution(ctx context.Context, storedEvent queries.Event, handler event.Handler[domain.Event], handlerExecutions map[string]queries.EventHandlerExecution) (queries.EventHandlerExecution, error) {
	var handlerExecution queries.EventHandlerExecution

	if handlerExecution, ok := handlerExecutions[handler.Name()]; ok {
		return handlerExecution, nil
	}

	handlerExecution, err := p.qs.CreateEventHandlerExecution(ctx, queries.CreateEventHandlerExecutionParams{
		ID:          mapping.PgUUID(uuid.NewString()),
		EventID:     storedEvent.ID,
		HandlerType: handler.Name(),
	})
	if err != nil {
		err = fmt.Errorf("failed to create new event handler execution: %w", err)
	}
	return handlerExecution, err
}
