package db

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
)

type Storage struct {
	q *Queries
}

func NewStorage(q *Queries) *Storage {
	return &Storage{
		q: q,
	}
}

func (s *Storage) Queries(ctx context.Context) *Queries {
	if tx, ok := transaction.FromContext(ctx); ok {
		return s.q.WithTx(tx)
	}
	return s.q
}
