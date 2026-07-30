package db

import (
	"context"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/db/transaction"
)

type Storage struct {
	q *sqlc.Queries
}

func NewStorage(q *sqlc.Queries) *Storage {
	return &Storage{
		q: q,
	}
}

func (s *Storage) Queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := transaction.FromContext(ctx); ok {
		return s.q.WithTx(tx)
	}
	return s.q
}
