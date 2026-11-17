package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/messenger/internal/media/domain/model/media"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/shared/app/apperr"
	"uuid"
)

var _ media.Repository = (*MediaRepository)(nil)

type MediaRepository struct {
	*repository.Repository
}

func NewMediaRepository(
	repo *repository.Repository,
) *MediaRepository {
	return &MediaRepository{
		Repository: repo,
	}
}

func (r *MediaRepository) Save(ctx context.Context, m *media.Media) error {
	err := r.Queries(ctx).UpsertMediaAsset(ctx, sqlc.UpsertMediaAssetParams{
		ID: m.ID(),
		Status: sqlc.AssetStatus(m.Status),
		Name: m.Name,
	})
	if err != nil {
		return fmt.Errorf("upsert media asset: %w", err)
	}

	r.RegisterAggregate(ctx, m)

	return nil
}

func (r *MediaRepository) Get(ctx context.Context, id uuid.UUID) (*media.Media, error) {
	a, err := r.Queries(ctx).GetMediaAsset(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get asset: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound("Media not found")
	}

	return media.Rehydrate(
		a.ID,
		a.Name,
		media.Status(a.Status),
	), nil
}
