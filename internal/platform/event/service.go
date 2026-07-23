package event

import (
	"bytes"
	"context"
	"encoding/json"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/platform/db"

	"github.com/google/uuid"
)

type Service struct {
	storage *db.Storage
}

func (s *EventStorage) StoreEvent(ctx context.Context, e domain.Event) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		return err
	}

	err := s.Queries(ctx).CreateEvent(ctx, db.CreateEventParams{
		ID:           mapping.PgUUID(uuid.New().String()),
		EventType:    e.Name(),
		EventPayload: buf.Bytes(),
	})

	return err
}

// func (s *EventStorage) ListUnprocessedEvents(ctx context.Context) ([]*StoredEvent, error) {
// 	events, err := s.Queries(ctx).ListUnprocessedEvents(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("event processor: failed to retreive events: %w", err)
// 	}

// 	storedEvents := make([]*StoredEvent, 0, len(events))
// 	for _, e := range events {
// 		storedEvents = append(storedEvents, &StoredEvent{
// 			ID:           e.ID.String(),
// 			EventType:    e.EventType,
// 			EventPayload: bytes.NewBuffer(e.EventPayload),
// 			OccurredAt:   e.OccurredAt.Time,
// 			ProcessedAt:  e.OccurredAt.Time,
// 		})
// 	}

// 	return storedEvents, nil
// }

// func (s *EventStorage) MarkEvenntAsProcessed(ctx context.Context, event *StoredEvent) ([]*StoredEvent, error) {
// 	updateParams := db.UpdateEventHandlerExecutionParams{
// 		ID:       execution.ID,
// 		Attempts: execution.Attempts + 1,
// 	}
// 	if err != nil {
// 		updateParams.Error = mapping.PgText(err.Error())
// 		updateParams.NextRetryAt = mapping.ToDBTimestamp(time.Now().Add(5 * time.Minute))
// 	} else {
// 		updateParams.Error = mapping.PgText("")
// 		updateParams.NextRetryAt = mapping.ToDBTimestamp(time.Time{})
// 	}
// 	if e := p.qs.UpdateEventHandlerExecution(ctx, updateParams); e != nil {
// 		err = fmt.Errorf("failed to update event handler execution: %w; %w", e, err)
// 	}
// 	return err
// }
