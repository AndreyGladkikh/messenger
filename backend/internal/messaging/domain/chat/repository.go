package chat

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, c *Chat) error
	AddPrivate(ctx context.Context, c *Chat, privateChatPair *PrivateChatPair) error
	Get(ctx context.Context, id uuid.UUID) (*Chat, error)
	// ListForUser(ctx context.Context, userID string) ([]*Chat, error)
	PrivateChatExists(ctx context.Context, participant1, participant2 uuid.UUID) (bool, error)
}
