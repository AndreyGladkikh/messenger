package chat_participant

import "context"

type Repository interface {
	Add(ctx context.Context, participants []*ChatParticipant) error
}