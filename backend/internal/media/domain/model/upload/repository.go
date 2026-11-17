package upload

import (
	"context"
	"uuid"
)

type Repository interface {
	Save(ctx context.Context, u *Upload) error
	Get(ctx context.Context, id uuid.UUID) (*Upload, error)
}
