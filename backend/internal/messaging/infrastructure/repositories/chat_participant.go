package repositories

import (
	"context"
	"messenger/messenger/internal/messaging/domain/chat_participant"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/db"
)

type ChatParticipantRepository struct {
	Repository
}

func NewChatParticipantRepository(
	q *sqlc.Queries,
) *ChatParticipantRepository {
	return &ChatParticipantRepository{
		Repository: Repository{
			q: q,
		},
	}
}

func (r *ChatParticipantRepository) Add(ctx context.Context, participants []*chat_participant.ChatParticipant) error {
	params := make([]sqlc.AddParticipantsToChatParams, 0, len(participants))
	for _, p := range participants {
		params = append(params, sqlc.AddParticipantsToChatParams{
			ID:            p.ID(),
			ChatID:        db.ToDBOptionalUUID(p.ChatID()),
			ParticipantID: p.ParticipantID(),
			Role:          string(p.Role()),
		})
		r.RegisterAggregate(ctx, p)
	}
	_, err := r.queries(ctx).AddParticipantsToChat(ctx, params)
	if err != nil {
		return err
	}

	return err
}
