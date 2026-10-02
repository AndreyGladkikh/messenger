package message

import (
	"context"

	"uuid"
)

type Repository interface {
	Add(ctx context.Context, m *Message) error
	ListForChat(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*Message, error)
}
