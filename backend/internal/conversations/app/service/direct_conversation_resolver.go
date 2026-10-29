package service

import (
	"context"
	"errors"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/shared/app/apperr"
	"uuid"
)

type DirectConversationResolver struct {
	conversationRepository conversation.Repository
	membershipRepository   membership.Repository
	participantLookup      ParticipantLookup
}

func NewDirectConversationResolver(
	conversationRepository conversation.Repository,
	membershipRepository membership.Repository,
	participantLookup ParticipantLookup,
) *DirectConversationResolver {
	return &DirectConversationResolver{
		conversationRepository: conversationRepository,
		membershipRepository:   membershipRepository,
		participantLookup:      participantLookup,
	}
}

func (d *DirectConversationResolver) FindOrCreate(
	ctx context.Context,
	participantOneID uuid.UUID,
	participantTwoID uuid.UUID,
) (*conversation.Conversation, bool, error) {
	participants, err := conversation.NewDirectConversationParticipantsPair(participantOneID, participantTwoID)
	if err != nil {
		return nil, false, err
	}

	cnv, err := d.conversationRepository.GetDirectBetween(ctx, participants)
	if err == nil {
		return cnv, true, nil
	}
	if !errors.Is(err, apperr.ErrNotFound) {
		return nil, false, err
	}

	cnv = conversation.New(conversation.ConversationTypeDirect)

	err = d.conversationRepository.SaveDirect(ctx, cnv, participants)
	if err != nil {
		return nil, false, err
	}

	for _, participantID := range participants.List() {
		exists, err := d.participantLookup.Exists(ctx, participantID)
		if err != nil {
			return nil, false, err
		}

		if !exists {
			return nil, false, apperr.NewNotFound("Participant doesn't exist")
		}

		membership := membership.New(
			cnv.ID(),
			participantID,
			membership.RoleMember,
		)
		err = d.membershipRepository.Save(ctx, membership)
		if err != nil {
			return nil, false, err
		}
	}

	return cnv, false, nil
}
