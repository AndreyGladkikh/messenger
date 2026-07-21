package repositories

import (
	"context"
	"messenger/messenger/internal/adapters/out/postgres/mapping"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/domain/chat_participant"

	"github.com/google/uuid"
)

type ChatParticipantRepository struct {
	Repository
}

func NewChatParticipantRepository(
	q *queries.Queries,
) *ChatRepository {
	return &ChatRepository{
		Repository: Repository{
			q: q,
		},
	}
}

func (r *ChatParticipantRepository) Add(ctx context.Context, participants []*chat_participant.ChatParticipant) error {
	params := make([]queries.AddParticipantsToChatParams, len(participants))
	for _, p := range participants {
		params = append(params, queries.AddParticipantsToChatParams{
			ID: mapping.PgUUID(uuid.NewString()),
			ChatID: mapping.PgUUID(p.ChatID()),
			ParticipantID: mapping.PgUUID(p.ParticipantID()),
			Role: string(p.Role()),
		})
	}
	_, err := r.queries(ctx).AddParticipantsToChat(ctx, params)
	if err != nil {
		return err
	}

	r.RegisterAggregate(ctx, participants...)

	return err
}