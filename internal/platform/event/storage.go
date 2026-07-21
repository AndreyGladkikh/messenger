package event

import (
	"bytes"
	"context"
	"encoding/json"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/domain/event"

	"github.com/google/uuid"
)

type EventStorage struct {
	q *queries.Queries
}

func NewEventStorage(
	q *queries.Queries,
) *EventStorage {
	return &EventStorage{
		q: q,
	}
}

func (r *EventStorage) queries(ctx context.Context) *queries.Queries {
	if tx, ok := transaction.FromContext(ctx); ok {
		return r.q.WithTx(tx)
	}
	return r.q
}

func (r *EventStorage) Add(ctx context.Context, e event.DomainEvent) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		return err
	}

	err := r.queries(ctx).CreateEvent(ctx, queries.CreateEventParams{
		ID: mapping.PgUUID(uuid.New().String()),
		EventType: e.Name(),
		EventPayload: buf.Bytes(),
	})

	return err
}
