package repositories

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/domain/chat"

	"github.com/google/uuid"
)

type ChatRepository struct {
	Repository
}

func (r *ChatRepository) Add(ctx context.Context, chat *chat.Chat) error {
	err := r.queries(ctx).CreateChat(ctx, queries.CreateChatParams{
		ID:   uuid.MustParse(chat.ID()),
		Type: string(chat.Type()),
		Name: mapping.OptionalString(chat.Name()),
	})
	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, chat)

	return err
}
