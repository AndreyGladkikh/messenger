package query_reader

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type Reader struct {
	*db.Store
}

func NewReader(
	store *db.Store,
) *Reader {
	return &Reader{
		Store: store,
	}
}

func (r *Reader) GetAccount(ctx context.Context, id uuid.UUID) (query.AccountView, error) {
	var zero query.AccountView

	acc, err := r.Queries(ctx).GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return zero, domain.ErrNotFound
		}
		return zero, err
	}

	return query.AccountView{
		ID:   acc.ID,
		Type: acc.Type,
	}, nil
}
