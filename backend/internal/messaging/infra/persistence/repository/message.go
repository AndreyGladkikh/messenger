package repository

import (
	"context"
	"messenger/messenger/internal/messaging/domain/model/message"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
)

type MessageRepository struct {
	*repository.Repository
}

func NewMessageRepository(repo *repository.Repository) *MessageRepository {
	return &MessageRepository{
		Repository: repo,
	}
}

func (m *MessageRepository) Save(ctx context.Context, msg *message.Message) error {
	err := m.Queries(ctx).InsertMessage(ctx, sqlc.InsertMessageParams{
		ID:             msg.ID(),
		SenderID:       msg.SenderID,
		ConversationID: msg.ConversationID,
		Body:           msg.Body,
		CreatedAt:      msg.CreatedAt,
	})
	if err != nil {
		return err
	}

	m.RegisterAggregate(ctx, msg)
	return nil
}
