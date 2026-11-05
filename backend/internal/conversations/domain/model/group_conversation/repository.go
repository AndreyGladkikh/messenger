package group_conversation

import (
	"context"
	"uuid"
)

type Repository interface {
	Save(context.Context, *GroupConversation) error
	Get(context.Context, uuid.UUID) (*GroupConversation, error)
}
