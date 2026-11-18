package repository

import (
	"context"
	"fmt"
	"messenger/messenger/internal/directory/domain/model"
	"messenger/messenger/internal/directory/domain/model/handle"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
)

var _ handle.Repository = (*Handle)(nil)

type Handle struct {
	*repository.Repository
}

func NewHandle(
	r *repository.Repository,
) *Handle {
	return &Handle{
		Repository: r,
	}
}

func (r *Handle) Save(ctx context.Context, h *handle.Handle) error {
	err := r.Queries(ctx).UpsertDirectoryHandle(ctx, sqlc.UpsertDirectoryHandleParams{
		TargetID: h.Target.ID,
		TargetType: h.Target.Type.String(),
		Handle: h.Value,
	})
	if err != nil {
		return fmt.Errorf("upsert directory entry: %w", err)
	}

	r.RegisterAggregate(ctx, h)

	return nil
}

func (r *Handle) DeleteByTarget(ctx context.Context, trgt model.TargetRef) error {
	err := r.Queries(ctx).DeleteDirectoryHandle(ctx, sqlc.DeleteDirectoryHandleParams{
		TargetID: trgt.ID,
		TargetType: trgt.Type.String(),
	})
	if err != nil {
		return fmt.Errorf("delete directory hadle")
	}

	return nil
}
