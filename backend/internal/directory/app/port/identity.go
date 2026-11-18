package port

import (
	"context"
	"uuid"
)

// todo rename UserExists
type IdentityReader interface {
	ExistsUser(ctx context.Context, userID uuid.UUID) (bool, error)
}
