package user_profile

import (
	"context"

	"uuid"
)

type Repository interface {
	Save(context.Context, *UserProfile) error
	Get(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
}
