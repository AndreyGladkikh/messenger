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
	id, err := uuid.Parse(chat.ID())
	if err != nil {
		return err
	}
	name := sql.NullString{String: chat.Name(), Valid: chat.Name() != ""}	

	err = r.qs.CreateChat(ctx, queries.CreateChatParams{
		ID: id,
		Type: string(chat.Type()),
		Name: name,
	})
	return err
}