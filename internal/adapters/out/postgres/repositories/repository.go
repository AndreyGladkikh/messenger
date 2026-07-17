package repositories

import (
	"context"
	"database/sql"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/adapters/out/postgres/transaction"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/platform/uow"
)

type Repository struct {
	db *sql.DB
	q *queries.Queries
}

func (r *Repository) queries(ctx context.Context) *queries.Queries {
	if tx, ok := transaction.FromContext(ctx); ok {
		return r.q.WithTx(tx)
	}
	return r.q
}

func (r *Repository) RegisterAggregate(ctx context.Context, aggregate domain.Aggregate) {
	if uow, ok := uow.FromContext(ctx); ok {
		uow.RegisterAggregate(aggregate)
	}
}