package chat

import "context"

type Repository interface {
	Add(ctx context.Context, c *Chat) error
	Get(ctx context.Context, id string) (*Chat, error)
	// ListForUser(ctx context.Context, userID string) ([]*Chat, error)
	PrivateChatExists(ctx context.Context, participant1, participant2 string) (bool, error)
}
