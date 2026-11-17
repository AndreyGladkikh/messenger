package port

import (
	"context"
	"uuid"
)

type MediaReader interface {
	GetAvatarURLs(ctx context.Context, avatarIDs []uuid.UUID) (map[uuid.UUID]string, error)
}
