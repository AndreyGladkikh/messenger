package repository

import (
	"context"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/platform/uow"
	sharedDomain "messenger/messenger/internal/shared/domain"
)

type Repository struct {
	*db.Storage
}

func NewRepository(storage *db.Storage) *Repository {
	return &Repository{
		Storage: storage,
	}
}

func (r *Repository) RegisterAggregate(ctx context.Context, aggregate sharedDomain.Aggregate) {
	if uow, ok := uow.FromContext(ctx); ok {
		uow.RegisterAggregate(aggregate)
	}
}
