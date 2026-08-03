package repository

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/auth/domain/user"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/postgres"
	"messenger/messenger/internal/platform/repository"
	sharedDomain "messenger/messenger/internal/shared/domain"

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

func (r *UserRepository) Add(ctx context.Context, u *user.User) error {
	err := r.Queries(ctx).CreateUser(ctx, sqlc.CreateUserParams{
		ID:           u.ID,
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		if err, ok := errors.AsType[*pgconn.PgError](err); ok && err.Code == postgres.UniqueViolationErrCode && err.ConstraintName == "users_login_key" {
			return user.ErrLoginAlreadyExists
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
		return nil, sharedDomain.ErrNotFound
	}

	return user.Rehydrate(
		u.ID,
		u.Login,
		u.PasswordHash,
		u.Name.String,
	), nil
}
