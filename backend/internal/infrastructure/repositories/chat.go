package repositories

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/domain/chat"
	"messenger/messenger/internal/infrastructure/db"
	"messenger/messenger/internal/infrastructure/sqlc"
)

type ChatRepository struct {
	Repository
}

func NewChatRepository(
	q *sqlc.Queries,
) *ChatRepository {
	return &ChatRepository{
		Repository: Repository{
			q: q,
		},
	}
}

func (r *ChatRepository) Add(ctx context.Context, chat *chat.Chat) error {
	err := r.queries(ctx).CreateChat(ctx, sqlc.CreateChatParams{
		ID:   db.ToDBUUID(chat.ID()),
		Type: string(chat.Type()),
		Name: db.ToDBText(chat.Name()),
	})
	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, chat)

	return err
}

func (r *ChatRepository) PrivateChatExists(ctx context.Context, participant1, participant2 string) (bool, error) {
	exists, err := r.queries(ctx).PrivateChatExists(ctx, sqlc.PrivateChatExistsParams{
		ParticipantID:   db.ToDBUUID(participant1),
		ParticipantID_2: db.ToDBUUID(participant2),
	})
	if err != nil {
		return false, err
	}

	return bool(exists), nil
}

func (r *ChatRepository) Get(ctx context.Context, chatID string) (*chat.Chat, error) {
	chatRow, err := r.queries(ctx).GetChat(ctx, db.ToDBUUID(chatID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, chat.ErrNotFound
	}

	chat := chat.Rehydrate(
		chatRow.ID.String(),
		chat.ChatType(chatRow.Type),
		chatRow.Name.String,
	)

	return chat, nil
}
