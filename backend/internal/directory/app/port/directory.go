package port

import (
	"context"
	"messenger/messenger/internal/directory/app/view"
	"uuid"
)

type DirectoryReader interface {
	GetUserProfilesByIDs(ctx context.Context, ids []uuid.UUID) ([]view.UserProfile, error)
}
