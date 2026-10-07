package reader

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type Reader struct {
	*store.Store
}

func NewReader(
	store *store.Store,
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
		Type: string(acc.Type),
	}, nil
}
