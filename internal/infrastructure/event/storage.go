package event

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/infrastructure/db"
	"time"

	"github.com/google/uuid"
)

var ErrHandlerAlreadyExecuted = errors.New("handler already executed")

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
		EventID: uuid.New().String(),
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
		ID:        event.ID,
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
	errs []error,
	nextRetryAt time.Time,
) error {
	errors := make([]string, 0, len(errs))
	for _, e := range errs {
		errors = append(errors, e.Error())
	}
	err := s.Queries(ctx).ProcessEvent(ctx, db.ProcessEventParams{
		ID:          event.ID,
		Status:      string(status),
		ClaimedAt:   mapping.ToDBTimestamp(time.Now()),
		Attempts:    event.Attempts + 1,
		Errors:       errors,
		NextRetryAt: mapping.ToDBTimestamp(nextRetryAt),
	})
	if err != nil {
		return fmt.Errorf("event processor: failed to mark event as processed: %w", err)
	}
	return nil
}

func (s *EventStorage) RegisterEventHandlerExecution(ctx context.Context, eventID, handler string) error {
	_, err := s.Queries(ctx).CreateInbox(ctx, db.CreateInboxParams{
		EventID: mapping.PgUUID(eventID),
		Handler: handler,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrHandlerAlreadyExecuted
	}
	return nil
}