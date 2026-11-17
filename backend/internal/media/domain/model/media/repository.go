package media

import (
	"context"
	"uuid"
)

type Repository interface {
	Save(ctx context.Context, m *Media) error
	Get(ctx context.Context, id uuid.UUID) (*Media, error)
}
