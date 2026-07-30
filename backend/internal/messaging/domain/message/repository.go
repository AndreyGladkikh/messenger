package message

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, m *Message) error
	ListForChat(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*Message, error)
}
