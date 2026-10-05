package service

import (
	"context"
	"errors"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type DirectConversationResolver struct {
	conversationRepository conversation.Repository
	membershipRepository   membership.Repository
}

func NewDirectConversationResolver(
	conversationRepository conversation.Repository,
	membershipRepository   membership.Repository,
) *DirectConversationResolver {
	return &DirectConversationResolver{
		conversationRepository: conversationRepository,
		membershipRepository: membershipRepository,
	}
}

func (f *DirectConversationResolver) FindOrCreate(
	ctx context.Context,
	participantOneID uuid.UUID,
	participantTwoID uuid.UUID,
) (*conversation.Conversation, error) {
	participants, err := conversation.NewDirectConversationParticipantsPair(participantOneID, participantTwoID)
	if err != nil {
		return nil, err
	}

	cnv, err := f.conversationRepository.GetDirectBetween(ctx, participants)
	if err == nil {
		return cnv, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	cnv = conversation.New(conversation.ConversationTypeDirect)

	err = f.conversationRepository.SaveDirect(ctx, cnv, participants)
	if err != nil {
		return nil, err
	}

	for _, participantID := range participants.List() {
		membership := membership.New(
			cnv.ID(),
			participantID,
			membership.RoleMember,
		)
		err = f.membershipRepository.Save(ctx, membership)
		if err != nil {
			return nil, err
		}
	}

	return cnv, nil
}
