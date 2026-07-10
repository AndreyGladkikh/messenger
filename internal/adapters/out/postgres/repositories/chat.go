package repositories

import (
	"context"
	"database/sql"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/domain/chat"

	"github.com/google/uuid"
)

type ChatRepository struct {
	db *sql.DB
	qs *queries.Queries
}

func (r *ChatRepository) Add(ctx context.Context, chat *chat.Chat) error {
	id, err := uuid.Parse(chat.ID)
	

	r.qs.CreateChat(ctx, queries.CreateChatParams{})
}