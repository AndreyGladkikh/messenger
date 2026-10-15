package persistence

import (
	"context"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
)

type EntryRepository struct {
	*repository.Repository
}

func NewEntryRepository(
	r *repository.Repository,
) *EntryRepository {
	return &EntryRepository{
		Repository: r,
	}
}

func (r *EntryRepository) Save(ctx context.Context, e *entry.Entry) error {
	err := r.Queries(ctx).UpsertDirectoryEntry(ctx, sqlc.UpsertDirectoryEntryParams{
		ID:          e.ID(),
		TargetID:    e.TargetID,
		TargetType:  string(e.TargetType),
		DisplayName: e.DisplayName,
		Handle:      e.Handle,
	})
	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, e)

	return nil
}
