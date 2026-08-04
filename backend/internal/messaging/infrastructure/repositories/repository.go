package repositories

import (
	"context"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/db/transaction"
	"messenger/messenger/internal/platform/uow"
	sharedDomain "messenger/messenger/internal/shared/domain"
)

type Repository struct {
	q *sqlc.Queries
}

func (r *Repository) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := transaction.FromContext(ctx); ok {
		return r.q.WithTx(tx)
	}
	return r.q
}

func (r *Repository) RegisterAggregate(ctx context.Context, aggregate sharedDomain.Aggregate) {
	if uow, ok := uow.FromContext(ctx); ok {
		uow.RegisterAggregate(aggregate)
	}
}
