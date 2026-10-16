package persistence

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type QueryService struct {
	*store.Store
}

func NewQueryService(
	store *store.Store,
) *QueryService {
	return &QueryService{
		Store: store,
	}
}

func (r *QueryService) GetAccount(ctx context.Context, id uuid.UUID) (query.AccountView, error) {
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
