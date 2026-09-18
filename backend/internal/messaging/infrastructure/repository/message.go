package repository

import (
	"context"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/platform/repository"

	"github.com/google/uuid"
)

type MessageRepository struct {
	*repository.Repository
}

func NewMessageRepository(
	repo *repository.Repository,
) *MessageRepository {
	return &MessageRepository{
		Repository: repo,
	}
}

func (r *MessageRepository) Add(ctx context.Context, m *message.Message) error {
	err := r.Queries(ctx).CreateMessage(ctx, sqlc.CreateMessageParams{
		ID:               m.ID(),
		SenderID:         m.SenderID(),
		ChatID:           m.ChatID(),
		Body:             m.Body(),
		ReplyToMessageID: db.ToDBOptionalUUID(m.ReplyToMessageID()),
	})

	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, m)

	return nil
}

func (r *MessageRepository) ListForChat(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*message.Message, error) {
	if limit == 0 {
		limit = 100
	}
	if offset == 0 {
		offset = 0
	}

	messageRows, err := r.Queries(ctx).ListMessagesForChat(ctx, sqlc.ListMessagesForChatParams{
		ChatID: chatID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*message.Message, 0, len(messageRows))
	for _, row := range messageRows {
		messages = append(messages, message.Rehydrate(
			row.ID,
			row.ChatID,
			row.SenderID,
			row.Body,
			row.ReplyToMessageID.UUID,
		))
	}

	return messages, nil
}
