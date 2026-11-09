package persistence

import (
	"context"
	"errors"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/shared/app/apperr"

	"github.com/jackc/pgx/v5/pgconn"
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
		AvatarUrl:   e.AvatarURL,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == postgres.UniqueViolationErrCode {
			switch pgErr.ConstraintName {
			case "entries_handle_key":
				return apperr.NewAlreadyExists(
					"Handle already exists",
					apperr.WithCode(apperr.CodeHandleAlreadyExists),
				)
			case "entries_target_id_target_type_key":
				return apperr.NewAlreadyExists(
					"Directory entry already exists",
					apperr.WithCode(apperr.CodeDirectoryEntryAlreadyExists),
				)
			}
		}
		return err
	}

	r.RegisterAggregate(ctx, e)

	return nil
}
