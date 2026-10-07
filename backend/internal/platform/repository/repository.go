package repository

import (
	"context"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/platform/uow"
	"messenger/messenger/internal/shared/domain"
)

type Repository struct {
	*store.Store
}

func NewRepository(store *store.Store) *Repository {
	return &Repository{
		Store: store,
	}
}

func (r *Repository) RegisterAggregate(ctx context.Context, aggregate domain.Aggregate) {
	if uow, ok := uow.FromContext(ctx); ok {
		uow.RegisterAggregate(aggregate)
	}
}
