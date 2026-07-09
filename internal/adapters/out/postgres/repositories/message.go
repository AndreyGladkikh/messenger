package repositories

import (
	"context"
	"database/sql"
	"messenger/messenger/internal/domain/message"
)

type MessageRepository struct {
	db *sql.DB
}

func (r *MessageRepository) Add(ctx context.Context, m *message.Message) error {
	return nil
}