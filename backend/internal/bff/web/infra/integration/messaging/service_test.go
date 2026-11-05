package messaging

import (
	"context"
	"errors"
	"messenger/messenger/internal/bff/web/application/port"
	cnvQry "messenger/messenger/internal/conversations/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceListUserConversations(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	conversationID := uuid.New()
	messageID := uuid.New()
	senderID := uuid.New()
	avatarURL := "https://example.com/avatar.png"
	title := "Alice"
	sentAt := time.Now()
	handler := appQuery.HandlerFunc[cnvQry.ListUserConversations, []cnvQry.ConversationSummary](
		func(_ context.Context, qry cnvQry.ListUserConversations) ([]cnvQry.ConversationSummary, error) {
			assert.Equal(t, userID, qry.UserID)

			return []cnvQry.ConversationSummary{{
				ID:        conversationID,
				Type:      port.ConversationTypeDirect,
				Title:     &title,
				AvatarURL: &avatarURL,
				LastMessage: cnvQry.LastMessage{
					ID:       messageID,
					Body:     "Hello",
					SenderID: senderID,
					SentAt:   sentAt,
				},
			}}, nil
		},
	)
	service := NewService(handler)

	conversations, err := service.ListUserConvesations(t.Context(), userID)

	require.NoError(t, err)
	assert.Equal(t, []port.ConversationSummary{{
		ID:        conversationID,
		Type:      port.ConversationTypeDirect,
		Title:     &title,
		AvatarURL: &avatarURL,
		LastMessage: port.LastMessage{
			ID:       messageID,
			Body:     "Hello",
			SenderID: senderID,
			SentAt:   sentAt,
		},
	}}, conversations)
}

func TestServiceListUserConversationsWrapsHandlerError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("read conversations")
	handler := appQuery.HandlerFunc[cnvQry.ListUserConversations, []cnvQry.ConversationSummary](
		func(context.Context, cnvQry.ListUserConversations) ([]cnvQry.ConversationSummary, error) {
			return nil, expectedErr
		},
	)
	service := NewService(handler)

	_, err := service.ListUserConvesations(t.Context(), uuid.New())

	require.ErrorIs(t, err, expectedErr)
	assert.ErrorContains(t, err, "list user conversations")
}
