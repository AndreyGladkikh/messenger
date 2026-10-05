package repository

import (
	"context"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	"uuid"
)

type MembershipRepository struct {
	*repository.Repository
}

func NewMembershipRepository(repo *repository.Repository) *MembershipRepository {
	return &MembershipRepository{
		Repository: repo,
	}
}

func (m *MembershipRepository) Exists(ctx context.Context, conversationID uuid.UUID, participantID uuid.UUID) (bool, error) {
	return m.Queries(ctx).ConversationMembershipExists(ctx, sqlc.ConversationMembershipExistsParams{
		ConversationID: conversationID,
		ParticipantID:  participantID,
	})
}

func (m *MembershipRepository) Save(ctx context.Context, membership *membership.Membership) error {
	err := m.Queries(ctx).CreateConversationMembership(ctx, sqlc.CreateConversationMembershipParams{
		ID:             membership.ID(),
		ConversationID: membership.ConversationID,
		ParticipantID:  membership.ParticipantID,
		Role:           string(membership.Role),
	})
	if err != nil {
		return err
	}

	m.RegisterAggregate(ctx, membership)
	return nil
}
