package chat

import "context"

type Repository interface {
	Add(context.Context, *Chat) error
	Get(context.Context, string) (*Chat, error)
	ListForUser(context.Context, string) ([]*Chat, error)
}
