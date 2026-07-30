package repository

import (
	"context"
	"messenger/messenger/internal/messaging/domain"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/platform/uow"
)

type Repository struct {
	*db.Storage
}

func NewRepository(storage *db.Storage) *Repository {
	return &Repository{
		Storage: storage,
	}
}

func (r *Repository) RegisterAggregate(ctx context.Context, aggregate domain.Aggregate) {
	if uow, ok := uow.FromContext(ctx); ok {
		uow.RegisterAggregate(aggregate)
	}
}
