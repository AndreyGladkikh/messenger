package event

import (
	"context"
	"messenger/messenger/internal/platform/store"
)

type Store struct {
	*store.Store
}

func NewStore(
	store *store.Store,
) *Store {
	return &Store{
		Store: store,
	}
}

func (s *Store) GetNextUnprocessedEvent(ctx context.Context) (Envelope, error) {
	outbox, err := s.Queries(ctx).GetEventToProcess(ctx)
	if err != nil {
		return Envelope{}, err
	}

	return Envelope{
		ID: outbox.EventID,
		Type: outbox.EventType,
		OccurredAt: outbox.OccurredAt,
		Payload: outbox.EventPayload,
	}, nil
}
