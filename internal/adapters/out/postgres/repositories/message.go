package repositories

import (
	"context"
	"database/sql"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/domain/message"
	"messenger/messenger/internal/platform/uow"
	"time"

	"github.com/google/uuid"
)

type MessageRepository struct {
	db *sql.DB
	qs *queries.Queries
}

func NewMessageRepository(
	db *sql.DB,
	qs *queries.Queries,
) *MessageRepository {
	return &MessageRepository{
		db: db,
		qs: qs,
	}
}

func (r *MessageRepository) Add(ctx context.Context, m *message.Message) error {
	messageID, err := uuid.Parse(m.ID())
	if err != nil {
		return err
	}
	senderID, err := uuid.Parse(m.SenderID())
	if err != nil {
		return err
	}
	chatID, err := uuid.Parse(m.ChatID())
	if err != nil {
		return err
	}
	replyToMessageID := uuid.NullUUID{}
	if m.ReplyToMessageID() != "" {
		replyID, err := uuid.Parse(m.ReplyToMessageID())
		if err != nil {
			return err
		}
		replyToMessageID.UUID = replyID
		replyToMessageID.Valid = true
	}

	err = r.qs.CreateMessage(ctx, queries.CreateMessageParams{
		ID:       messageID,
		SenderID: senderID,
		ChatID:   chatID,
		Body:     m.Body(),
		ReplyToMessageID: replyToMessageID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if uo, ok := uow.FromContext(ctx); ok {
		uo.RegisterAggregate(m)
	}

	return err
}

func (r *MessageRepository) ListForChat(ctx context.Context, chatID string, limit, offset int) ([]*message.Message, error) {
	if limit == 0 {
		limit = 100
	}
	if offset == 0 {
		offset = 0
	}

	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		return nil, err
	}

	messageRows, err := r.qs.ListMessagesForChat(ctx, queries.ListMessagesForChatParams{
		ChatID: chatUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*message.Message, len(messageRows))
	for i, row := range messageRows {
		messages[i] = message.Rehydrate(
			row.ID.String(),
			row.ChatID.String(),
			row.SenderID.String(),
			row.Body,
			row.ReplyToMessageID.UUID.String(),
		)
	}

	return messages, nil
}
