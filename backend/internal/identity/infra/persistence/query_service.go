package persistence

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/shared/app/apperr"
	"uuid"
)

var _ query.Reader = (*QueryService)(nil)

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

func (q *QueryService) GetAccount(ctx context.Context, id uuid.UUID) (query.AccountView, error) {
	var zero query.AccountView

	acc, err := q.Queries(ctx).GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return zero, apperr.NewNotFound("Account not found by id")
		}
		return zero, err
	}

	return query.AccountView{
		ID:   acc.ID,
		Type: string(acc.Type),
	}, nil
}

func (q *QueryService) GetUserIDByLogin(ctx context.Context, login string) (uuid.UUID, error) {
	u, err := q.Queries(ctx).GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.UUID{}, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.UUID{}, apperr.NewNotFound("Not found user by login")
	}

	return u.ID, err
}
