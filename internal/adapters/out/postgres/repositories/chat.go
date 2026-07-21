package repositories

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/domain/chat"
)

type ChatRepository struct {
	Repository
}

func NewChatRepository(
	q *queries.Queries,
) *ChatRepository {
	return &ChatRepository{
		Repository: Repository{
			q: q,
		},
	}
}

func (r *ChatRepository) Add(ctx context.Context, chat *chat.Chat) error {
	err := r.queries(ctx).CreateChat(ctx, queries.CreateChatParams{
		ID:   mapping.PgUUID(chat.ID()),
		Type: string(chat.Type()),
		Name: mapping.PgText(chat.Name()),
	})
	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, chat)

	return err
}

func (r *ChatRepository) PrivateChatExists(ctx context.Context, participant1, participant2 string) (bool, error) {
	exists, err := r.queries(ctx).PrivateChatExists(ctx, queries.PrivateChatExistsParams{
		ParticipantID: mapping.PgUUID(participant1),
		ParticipantID_2: mapping.PgUUID(participant2),
	})
	if err != nil {
		return false, err
	}

	return bool(exists), nil
}