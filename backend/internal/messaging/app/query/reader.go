package query

import (
	"context"
	"uuid"
)

type Reader interface {
	ListConversationsForUser(ctx context.Context, userID uuid.UUID) ([]ConversationSummary, error)
}
