package event

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/shared/app/event"
	"time"

	"uuid"
)

var ErrHandlerAlreadyExecuted = errors.New("handler already executed")

type EventService struct {
	*store.Store
}

func NewEventService(
	store *store.Store,
) *EventService {
	return &EventService{
		Store: store,
	}
}

func (s *EventService) PutEventToOutbox(ctx context.Context, kind string, e event.Event) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		return err
	}

	err := s.Queries(ctx).PutToOutbox(ctx, sqlc.PutToOutboxParams{
		ID:           uuid.New(),
		Kind: sqlc.OutboxKind(kind),
		Type: e.Name(),
		Payload: buf.Bytes(),
	})

	return err
}

func (s *EventService) GetNextUnprocessedEvent(ctx context.Context) (sqlc.Outbox, error) {
	return s.Queries(ctx).GetEventToProcess(ctx)
}

func (s *EventService) UpdateEvent(
	ctx context.Context,
	event sqlc.Outbox,
	status Status,
	errCause error,
	nextRetryAt time.Time,
) error {
	errMsg := errCause.Error()
	err := s.Queries(ctx).ProcessEvent(ctx, sqlc.ProcessEventParams{
		ID:          event.ID,
		Status:      sqlc.OutboxStatus(status),
		Error: &errMsg,
		Attempts:    event.Attempts + 1,
		NextRetryAt: &nextRetryAt,
	})
	if err != nil {
		return fmt.Errorf("event processor: failed to mark event as processed: %w", err)
	}
	return nil
}

func (s *EventService) RegisterEventHandlerExecution(ctx context.Context, eventID uuid.UUID, handler string) error {
	_, err := s.Queries(ctx).CreateInbox(ctx, sqlc.CreateInboxParams{
		EventID: eventID,
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
