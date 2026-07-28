package db

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/infrastructure/sqlc"
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
