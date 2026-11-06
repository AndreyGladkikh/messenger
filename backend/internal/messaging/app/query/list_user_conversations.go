package query

import (
	"context"
	"fmt"
	"messenger/messenger/internal/shared/app/query"
	"time"
	"uuid"
)

const ListUserConversationsQueryName = "ListUserConversations"

var _ query.Query = ListUserConversations{}

type ListUserConversations struct {
	UserID uuid.UUID
}

func (l ListUserConversations) Name() string {
	return ListUserConversationsQueryName
}

type ConversationSummary struct {
	ID          uuid.UUID
	Type        string
	Title       *string
	AvatarURL   *string
	LastMessage LastMessage
	DirectParticipantID *uuid.UUID
}

type LastMessage struct {
	ID       uuid.UUID
	Body     string
	SenderID uuid.UUID
	SentAt   time.Time
}

type ListUserConversationsHandler struct {
	reader Reader
}

func NewListUserConversationsHandler(
	reader Reader,
) *ListUserConversationsHandler{
	return &ListUserConversationsHandler{
		reader: reader,
	}
}

func (h *ListUserConversationsHandler) Handle(ctx context.Context, qry ListUserConversations) ([]ConversationSummary, error) {
	cnvs, err := h.reader.ListConversationsForUser(ctx, qry.UserID)
	if err != nil {
		return nil, fmt.Errorf("list conversations for user: %w", err)
	}

	return cnvs, nil
}
