package conversation

import (
	"context"
	"uuid"
)

type Repository interface {
	Save(context.Context, *Conversation) error
	Get(context.Context, uuid.UUID) (*Conversation, error)
	SaveDirect(context.Context, *Conversation, DirectConversationParticipantsPair) error
	GetDirectBetween(ctx context.Context, participants DirectConversationParticipantsPair) (*Conversation, error)
}
