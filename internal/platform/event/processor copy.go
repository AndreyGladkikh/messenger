package event

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"messenger/messenger/internal/adapters/out/logger"
// 	"messenger/messenger/internal/adapters/out/postgres/mapping"
// 	"messenger/messenger/internal/adapters/out/postgres/transaction"
// 	"messenger/messenger/internal/application/event"
// 	"messenger/messenger/internal/domain"
// 	"messenger/messenger/internal/platform/db"
// 	"time"

// 	"github.com/google/uuid"
// )

// type ProcessorConfig struct {
// }

// type Processor struct {
// 	eventStorage *EventStorage
// 	txManager    *transaction.Manager
// 	logger       *logger.Logger
// 	qs           *db.Queries
// 	eventBus     *Bus
// 	cfg          *ProcessorConfig
// 	inShutdown   bool
// }

// func NewProcessor(
// 	eventStorage *EventStorage,
// 	txManager *transaction.Manager,
// 	logger *logger.Logger,
// 	qs *db.Queries,
// 	eventBus *Bus,
// ) *Processor {
// 	return &Processor{
// 		eventStorage: eventStorage,
// 		txManager:    txManager,
// 		logger:       logger,
// 		qs:           qs,
// 		eventBus:     eventBus,
// 	}
// }

// func (p *Processor) Run(ctx context.Context) error {
// 	go func() {
// 		<-ctx.Done()
// 		p.Shutdown()
// 	}()

// 	for {
// 		if p.inShutdown {
// 			return nil
// 		}

// 		err := p.txManager.WithTransaction(ctx, func(ctx context.Context) error {
// 			event, err := p.eventStorage.GetNextUnprocessedEvent(ctx)
// 			if err != nil && !errors.Is(err, sql.ErrNoRows) {
// 				return fmt.Errorf("retreive event to process: %w", err)
// 			}
// 			if errors.Is(err, sql.ErrNoRows) {
// 				time.Sleep(1*time.Second)
// 				return nil
// 			}

// 			err = p.processEvent(ctx, event)
// 			if err != nil {
// 				return fmt.Errorf("failed to process event: %w", err)
// 			}

// 			err = p.eventStorage.MarkEventAsProcessed(ctx, event)
// 			if err != nil {
// 				return err
// 			}

// 			return nil
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		// events, err := p.qs.ListUnprocessedEvents(ctx)
// 		// if err != nil {
// 		// 	return fmt.Errorf("event processor: failed to retreive events: %w", err)
// 		// }
// 		// if len(events) == 0 {
// 		// 	time.Sleep(1*time.Second)
// 		// 	continue
// 		// }

// 		// for _, e := range events {
// 		// 	go func() {
// 		// 		err := p.processEvent(ctx, e)
// 		// 		if err != nil {
// 		// 			p.logger.Error(
// 		// 				"event processor: failed to process event",
// 		// 				"eventID", e.ID,
// 		// 				"error", err.Error(),
// 		// 			)
// 		// 		}
// 		// 	}()
// 		// 	err := p.qs.ProcessEvent(ctx, db.ProcessEventParams{
// 		// 		ID: e.ID,
// 		// 		ProcessedAt: mapping.ToDBTimestamp(time.Now()),
// 		// 	})
// 		// 	if err != nil {
// 		// 		return fmt.Errorf("event processor: failed to mark event as processed: %w", err)
// 		// 	}
// 		// }
// 	}
// }

// func (p *Processor) Shutdown() error {
// 	p.inShutdown = true
// 	return nil
// }

// func (p *Processor) processEvent(ctx context.Context, storedEvent db.Event) error {
// 	event, err := translateStoredEventToDispatchedEvent(storedEvent)
// 	if err != nil {
// 		return fmt.Errorf("event processor: failed to translate stored event to dispatched event")
// 	}

// 	handlers := p.eventBus.HandlersForEvent(event)
// 	handlerExecutions, err := p.eventStorage.GetEventHandlerExecutionsByEventId(ctx, storedEvent.ID.String())
// 	for _, h := range handlers {
// 		if !shouldExecuteHandler(h, handlerExecutions) {
// 			continue
// 		}

// 		err := p.executeHandler(ctx, storedEvent, event, h, handlerExecutions)
// 		if err != nil {
// 			return fmt.Errorf("event processor: failed to execute handler %q: %w", h.Name(), err)
// 		}
// 	}

// 	return nil
// }

// func (p *Processor) executeHandler(ctx context.Context, storedEvent db.Event, event domain.Event, handler event.Handler[domain.Event], handlerExecutions map[string]db.EventHandlerExecution) error {
// 	execution, err := p.getExecution(ctx, storedEvent, handler, handlerExecutions)
// 	if err != nil {
// 		return err
// 	}

// 	err = p.eventBus.DispatchForHandler(ctx, event, handler)
// 	if err != nil {
// 		err = fmt.Errorf("failed to execute handler %q for event %q: %w", handler.Name(), event.Name(), err)
// 	}

// 	return p.eventStorage.UpdateExecution(ctx, execution, err)
// }

// func shouldExecuteHandler[T domain.Event](handler event.Handler[T], handlerExecutions map[string]db.EventHandlerExecution) bool {
// 	var execution db.EventHandlerExecution

// 	execution, ok := handlerExecutions[handler.Name()]
// 	if !ok {
// 		return true
// 	}

// 	if execution.NextRetryAt.Valid && time.Now().Before(execution.NextRetryAt.Time) {
// 		return false
// 	}
// 	if execution.Attempts > 0 && !execution.Error.Valid {
// 		return false
// 	}
// 	if execution.Attempts >= 3 {
// 		return false
// 	}

// 	return true
// }

// func (p *Processor) getExecution(ctx context.Context, storedEvent db.Event, handler event.Handler[domain.Event], handlerExecutions map[string]db.EventHandlerExecution) (db.EventHandlerExecution, error) {
// 	var handlerExecution db.EventHandlerExecution

// 	if handlerExecution, ok := handlerExecutions[handler.Name()]; ok {
// 		return handlerExecution, nil
// 	}

// 	handlerExecution, err := p.qs.CreateEventHandlerExecution(ctx, db.CreateEventHandlerExecutionParams{
// 		ID:          mapping.PgUUID(uuid.NewString()),
// 		EventID:     storedEvent.ID,
// 		HandlerType: handler.Name(),
// 	})
// 	if err != nil {
// 		err = fmt.Errorf("failed to create new event handler execution: %w", err)
// 	}
// 	return handlerExecution, err
// }

// // func (p *Processor) processEvent(ctx context.Context, storedEvent db.Event) error {
// // 	event, err := translateStoredEventToDispatchedEvent(storedEvent)
// // 	if err != nil {
// // 		return fmt.Errorf("event processor: failed to translate stored event to dispatched event")
// // 	}
// // 	handlers := p.eventBus.HandlersForEvent(event)
// // 	handlerExecutions, err := p.qs.GetEventHandlerExecutionsByEventId(ctx, storedEvent.ID)
// // 	if err != nil {
// // 		return fmt.Errorf("event processor: failed to retrieve event handler executions: %w", err)
// // 	}
// // 	handlerExecutionsMap := make(map[string]db.EventHandlerExecution)
// // 	for _, execution := range handlerExecutions {
// // 		handlerExecutionsMap[execution.HandlerType] = execution
// // 	}

// // 	for _, h := range handlers {
// // 		if !shouldExecuteHandler(h, handlerExecutionsMap) {
// // 			continue
// // 		}

// // 		err := p.executeHandler(ctx, storedEvent, event, h, handlerExecutionsMap)
// // 		if err != nil {
// // 			return fmt.Errorf("event processor: failed to execute handler %q: %w", h.Name(), err)
// // 		}
// // 	}

// // 	return nil
// // }

// // func (p *Processor) executeHandler(ctx context.Context, storedEvent db.Event, event domain.Event, handler event.Handler[domain.Event], handlerExecutions map[string]db.EventHandlerExecution) error {
// // 	execution, err := p.getExecution(ctx, storedEvent, handler, handlerExecutions)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	err = p.eventBus.DispatchForHandler(ctx, event, handler)
// // 	if err != nil {
// // 		err = fmt.Errorf("failed to execute event handler: %w", err)
// // 	}

// // 	updateParams := db.UpdateEventHandlerExecutionParams{
// // 		ID:       execution.ID,
// // 		Attempts: execution.Attempts + 1,
// // 	}
// // 	if err != nil {
// // 		updateParams.Error = mapping.PgText(err.Error())
// // 		updateParams.NextRetryAt = mapping.ToDBTimestamp(time.Now().Add(5 * time.Minute))
// // 	} else {
// // 		updateParams.Error = mapping.PgText("")
// // 		updateParams.NextRetryAt = mapping.ToDBTimestamp(time.Time{})
// // 	}
// // 	if e := p.qs.UpdateEventHandlerExecution(ctx, updateParams); e != nil {
// // 		err = fmt.Errorf("failed to update event handler execution: %w; %w", e, err)
// // 	}
// // 	return err
// // }

// // func shouldExecuteHandler[T domain.Event](handler event.Handler[T], handlerExecutions map[string]db.EventHandlerExecution) bool {
// // 	var execution db.EventHandlerExecution

// // 	execution, ok := handlerExecutions[handler.Name()]
// // 	if !ok {
// // 		return true
// // 	}

// // 	if execution.NextRetryAt.Valid && time.Now().Before(execution.NextRetryAt.Time) {
// // 		return false
// // 	}
// // 	if execution.Attempts > 0 && !execution.Error.Valid {
// // 		return false
// // 	}
// // 	if execution.Attempts >= 3 {
// // 		return false
// // 	}

// // 	return true
// // }

// // func (p *Processor) getExecution(ctx context.Context, storedEvent db.Event, handler event.Handler[domain.Event], handlerExecutions map[string]db.EventHandlerExecution) (db.EventHandlerExecution, error) {
// // 	var handlerExecution db.EventHandlerExecution

// // 	if handlerExecution, ok := handlerExecutions[handler.Name()]; ok {
// // 		return handlerExecution, nil
// // 	}

// // 	handlerExecution, err := p.qs.CreateEventHandlerExecution(ctx, db.CreateEventHandlerExecutionParams{
// // 		ID:          mapping.PgUUID(uuid.NewString()),
// // 		EventID:     storedEvent.ID,
// // 		HandlerType: handler.Name(),
// // 	})
// // 	if err != nil {
// // 		err = fmt.Errorf("failed to create new event handler execution: %w", err)
// // 	}
// // 	return handlerExecution, err
// // }
