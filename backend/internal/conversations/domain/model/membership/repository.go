package membership

import (
	"context"
	"uuid"
)

type Repository interface {
	Exists(ctx context.Context, conversationID uuid.UUID, participantID uuid.UUID) (bool, error)
	Save(ctx context.Context, membership *Membership) error
}
