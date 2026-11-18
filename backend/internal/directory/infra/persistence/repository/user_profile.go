package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/messenger/internal/directory/domain/model/user_profile"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/shared/app/apperr"
	"uuid"
)

var _ user_profile.Repository = (*UserProfile)(nil)

type UserProfile struct {
	*repository.Repository
}

func NewUserProfile(
	r *repository.Repository,
) *UserProfile {
	return &UserProfile{
		Repository: r,
	}
}

func (r *UserProfile) Save(ctx context.Context, up *user_profile.UserProfile) error {
	err := r.Queries(ctx).UpsertUserProfile(ctx, sqlc.UpsertUserProfileParams{
		UserID:      up.UserID,
		DisplayName: up.DisplayName,
		AvatarID:    up.AvatarID,
		Bio:         up.Bio,
		Birthday:    up.Birthday,
	})
	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, up)

	return nil
}

func (r *UserProfile) Get(ctx context.Context, userID uuid.UUID) (*user_profile.UserProfile, error) {
	up, err := r.Queries(ctx).GetUserProfile(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound(fmt.Sprintf("User profile for user with id %q not found", userID))
	}

	p := user_profile.Rehydrate(
		up.UserID,
		up.DisplayName,
		up.AvatarID,
		up.Bio,
		up.Birthday,
	)
	return p, nil
}
