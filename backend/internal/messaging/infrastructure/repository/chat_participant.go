package repository

import (
	"context"
	"messenger/messenger/internal/messaging/domain/chat_participant"
	"messenger/messenger/internal/messaging/infrastructure/sqlc"
	"messenger/messenger/internal/platform/db"
	"messenger/messenger/internal/shared/infrastructure/repository"
)

type ChatParticipantRepository struct {
	*repository.Repository
}

func NewChatParticipantRepository(
	repo *repository.Repository,
) *ChatParticipantRepository {
	return &ChatParticipantRepository{
		Repository: repo,
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
	_, err := r.Queries(ctx).AddParticipantsToChat(ctx, params)
	if err != nil {
		return err
	}

	return err
}
