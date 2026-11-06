package port

import (
	"context"
	"time"
	"uuid"
)

const (
	ConversationTypeDirect = "direct"
	ConversationTypeGroup  = "group"
)

type ConversationSummary struct {
	ID          uuid.UUID
	Type        string
	Title       *string
	AvatarURL   *string
	LastMessage LastMessage
}

type LastMessage struct {
	ID       uuid.UUID
	Body     string
	SenderID uuid.UUID
	SentAt   time.Time
}

type MessagingReader interface {
	ListUserConvesations(ctx context.Context, userID uuid.UUID) ([]ConversationSummary, error)
}
