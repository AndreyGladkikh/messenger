package membership

import (
	"context"
	"messenger/messenger/internal/conversation/domain/model/participant"
	"uuid"
)

type Repository interface {
	Exists(ctx context.Context, participant participant.Participant, conversationID uuid.UUID) (bool, error)
	Save(ctx context.Context, membership *Membership) error
}
