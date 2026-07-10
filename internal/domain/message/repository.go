package message

import "context"

type Repository interface {
	Add(ctx context.Context, m *Message) error
	ListForChat(ctx context.Context, chatID string) ([]*Message, error)
}