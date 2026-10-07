package db

import (
	"context"
	"messenger/messenger/internal/platform/db/transaction"
	"messenger/messenger/internal/platform/sqlc"
)

type Store struct {
	q *sqlc.Queries
}

func NewStore(q *sqlc.Queries) *Store {
	return &Store{
		q: q,
	}
}

func (s *Store) Queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := transaction.FromContext(ctx); ok {
		return s.q.WithTx(tx)
	}
	return s.q
}
