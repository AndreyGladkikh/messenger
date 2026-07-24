package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/platform/db"
	"time"

	"github.com/google/uuid"
)

type EventStorage struct {
	*db.Storage
}

func NewEventStorage(
	storage *db.Storage,
) *EventStorage {
	return &EventStorage{
		Storage: storage,
	}
}

func (s *EventStorage) Add(ctx context.Context, e domain.Event) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		return err
	}

	err := s.Queries(ctx).PutToOutbox(ctx, db.PutToOutboxParams{
		ID:           mapping.PgUUID(uuid.New().String()),
		EventType:    e.Name(),
		EventPayload: buf.Bytes(),
	})

	return err
}

func (s *EventStorage) GetNextUnprocessedEvent(ctx context.Context) (db.Outbox, error) {
	return s.Queries(ctx).GetEventToProcess(ctx)
}

func (s *EventStorage) MarkEventAsProcessed(ctx context.Context, event db.Outbox) error {
	err := s.Queries(ctx).ProcessEvent(ctx, db.ProcessEventParams{
		ID: event.ID,
		ClaimedAt: mapping.ToDBTimestamp(time.Now()),
	})
	if err != nil {
		return fmt.Errorf("event processor: failed to mark event as processed: %w", err)
	}
	return nil
}

func (s *EventStorage) UpdateEvent(
	ctx context.Context,
	event db.Outbox,
	status status,
	er string,
	nextRetryAt time.Time,
) error {
	err := s.Queries(ctx).ProcessEvent(ctx, db.ProcessEventParams{
		ID: event.ID,
		Status: string(status),
		ClaimedAt: mapping.ToDBTimestamp(time.Now()),
		Attempts: event.Attempts + 1,
		Error: mapping.PgText(er),
		NextRetryAt: mapping.ToDBTimestamp(nextRetryAt),
	})
	if err != nil {
		return fmt.Errorf("event processor: failed to mark event as processed: %w", err)
	}
	return nil
}

// func (s *EventStorage) GetEventHandlerExecutionsByEventId(ctx context.Context, eventID string) (map[string]db.EventHandlerExecution, error) {
// 	handlerExecutions, err := s.Queries(ctx).GetEventHandlerExecutionsByEventId(ctx, mapping.PgUUID(eventID))
// 	if err != nil {
// 		return nil, fmt.Errorf("event processor: failed to retrieve event handler executions: %w", err)
// 	}

// 	handlerExecutionsMap := make(map[string]db.EventHandlerExecution, len(handlerExecutions))
// 	for _, execution := range handlerExecutions {
// 		handlerExecutionsMap[execution.HandlerType] = execution
// 	}
// 	return handlerExecutionsMap, nil
// }

// func (s *EventStorage) UpdateExecution(ctx context.Context, execution db.EventHandlerExecution, execErr error) error {
// 	updateParams := db.UpdateEventHandlerExecutionParams{
// 		ID:       execution.ID,
// 		Attempts: execution.Attempts + 1,
// 	}
// 	if execErr != nil {
// 		updateParams.Error = mapping.PgText(execErr.Error())
// 		updateParams.NextRetryAt = mapping.ToDBTimestamp(time.Now().Add(5 * time.Minute))
// 	} else {
// 		updateParams.Error = mapping.PgText("")
// 		updateParams.NextRetryAt = mapping.ToDBTimestamp(time.Time{})
// 	}
// 	if err := s.Queries(ctx).UpdateEventHandlerExecution(ctx, updateParams); err != nil {
// 		return fmt.Errorf("update event handler execution: %w", err)
// 	}
// 	return nil
// }