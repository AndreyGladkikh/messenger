package message

import "context"

type Repository interface {
	NextID() string
	Add(ctx context.Context, m *Message) error
}