package repository

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/auth/domain/user"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/repository"
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
		ID: u.ID,
		Login: u.Login,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
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
		return nil, domain.ErrNotFound
	}

	return user.Rehydrate(
		u.ID,
		u.Login,
		u.PasswordHash,
		u.Name.String,
	), nil
}