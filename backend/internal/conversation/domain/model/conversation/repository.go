package conversation

import (
	"context"
	"messenger/messenger/internal/conversation/domain/model/participant"
	"uuid"
)

type Repository interface {
	Save(context.Context, *Conversation) error
	ExistsDirectBetween(ctx context.Context, participants DirectConversationParticipantsPair) (bool, error)
	GetDirectBetween(ctx context.Context, participants DirectConversationParticipantsPair) (*Conversation, error)
	ExistsDirectWith(ctx context.Context, conversationID uuid.UUID, participant participant.Participant) (bool, error)
}
