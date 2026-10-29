package repository

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/identity/app"
	"messenger/messenger/internal/identity/domain/model/user"
	"messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/shared/app/apperr"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	*repository.Repository
}

func NewUserRepository(repo *repository.Repository) *UserRepository {
	return &UserRepository{
		Repository: repo,
	}
}

func (r *UserRepository) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	return r.Queries(ctx).UserExistsByLogin(ctx, login)
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	err := r.Queries(ctx).CreateAccount(ctx, sqlc.CreateAccountParams{
		ID:   u.ID(),
		Type: sqlc.AccountTypeUser,
	})
	if err != nil {
		return err
	}

	err = r.Queries(ctx).CreateUser(ctx, sqlc.CreateUserParams{
		ID:           u.ID(),
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		if err, ok := errors.AsType[*pgconn.PgError](err); ok && err.Code == postgres.UniqueViolationErrCode && err.ConstraintName == "users_login_key" {
			return app.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	u, err := r.Queries(ctx).GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound("User not found by login")
	}

	return user.Rehydrate(
		u.ID,
		u.Login,
		u.PasswordHash,
	), nil
}
