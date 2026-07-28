package repositories

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/infrastructure/sqlc"
)

type MessageRepository struct {
	Repository
}

func NewMessageRepository(
	q *sqlc.Queries,
) *MessageRepository {
	return &MessageRepository{
		Repository: Repository{
			q: q,
		},
	}
}

func (r *MessageRepository) Add(ctx context.Context, m *message.Message) error {
	err := r.queries(ctx).CreateMessage(ctx, sqlc.CreateMessageParams{
		ID:               mapping.PgUUID(m.ID()),
		SenderID:         mapping.PgUUID(m.SenderID()),
		ChatID:           mapping.PgUUID(m.ChatID()),
		Body:             m.Body(),
		ReplyToMessageID: mapping.PgUUID(m.ReplyToMessageID()),
	})

	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, m)

	return nil
}

func (r *MessageRepository) ListForChat(ctx context.Context, chatID string, limit, offset int) ([]*message.Message, error) {
	if limit == 0 {
		limit = 100
	}
	if offset == 0 {
		offset = 0
	}

	messageRows, err := r.queries(ctx).ListMessagesForChat(ctx, sqlc.ListMessagesForChatParams{
		ChatID: mapping.PgUUID(chatID),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*message.Message, 0, len(messageRows))
	for _, row := range messageRows {
		messages = append(messages, message.Rehydrate(
			row.ID.String(),
			row.ChatID.String(),
			row.SenderID.String(),
			row.Body,
			row.ReplyToMessageID.String(),
		))
	}

	return messages, nil
}
