package port

import (
	"context"
	"uuid"
)

type MediaReader interface {
	MediaExists(ctx context.Context, mediaID uuid.UUID) (bool, error)
}
