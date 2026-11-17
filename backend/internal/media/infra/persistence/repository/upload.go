package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/messenger/internal/media/domain/model/upload"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/shared/app/apperr"
	"uuid"
)

var _ upload.Repository = (*UploadRepository)(nil)

type UploadRepository struct {
	*repository.Repository
}

func NewUploadRepository(
	repo *repository.Repository,
) *UploadRepository {
	return &UploadRepository{
		Repository: repo,
	}
}

func (r *UploadRepository) Save(ctx context.Context, u *upload.Upload) error {
	err := r.Queries(ctx).UpsertMediaUpload(ctx, sqlc.UpsertMediaUploadParams{
		ID: u.ID(),
		Status: sqlc.UploadStatus(u.Status),
		InitiatorID: u.InitiatorID,
		AssetID: u.MediaID,
	})
	if err != nil {
		return fmt.Errorf("upsert media upload: %w", err)
	}

	r.RegisterAggregate(ctx, u)

	return nil
}

func (r *UploadRepository) Get(ctx context.Context, id uuid.UUID) (*upload.Upload, error) {
	u, err := r.Queries(ctx).GetUpload(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get upload: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound("Upload not found")
	}

	return upload.Rehydrate(
		u.ID,
		u.InitiatorID,
		upload.Status(u.Status),
		u.AssetID,
	), nil
}
