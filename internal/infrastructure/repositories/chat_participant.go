package repositories

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/domain/chat_participant"
	"messenger/messenger/internal/infrastructure/db"
)

type ChatParticipantRepository struct {
	Repository
}

func NewChatParticipantRepository(
	q *db.Queries,
) *ChatParticipantRepository {
	return &ChatParticipantRepository{
		Repository: Repository{
			q: q,
		},
	}
}

func (r *ChatParticipantRepository) Add(ctx context.Context, participants []*chat_participant.ChatParticipant) error {
	params := make([]db.AddParticipantsToChatParams, len(participants))
	for i, p := range participants {
		params[i] = db.AddParticipantsToChatParams{
			ID:            mapping.PgUUID(p.ID()),
			ChatID:        mapping.PgUUID(p.ChatID()),
			ParticipantID: mapping.PgUUID(p.ParticipantID()),
			Role:          string(p.Role()),
		}
		r.RegisterAggregate(ctx, p)
	}
	_, err := r.queries(ctx).AddParticipantsToChat(ctx, params)
	if err != nil {
		return err
	}

	return err
}
