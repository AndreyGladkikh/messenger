package service

import (
	"context"
	"errors"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/shared/app/apperr"
	"uuid"
)

type DirectConversationResolver struct {
	conversationRepository conversation.Repository
	membershipRepository   membership.Repository
	participantLookup ParticipantLookup
}

func NewDirectConversationResolver(
	conversationRepository conversation.Repository,
	membershipRepository membership.Repository,
	participantLookup ParticipantLookup,
) *DirectConversationResolver {
	return &DirectConversationResolver{
		conversationRepository: conversationRepository,
		membershipRepository:   membershipRepository,
		participantLookup: participantLookup,
	}
}

func (d *DirectConversationResolver) FindOrCreate(
	ctx context.Context,
	participantOneID uuid.UUID,
	participantTwoID uuid.UUID,
) (*conversation.Conversation, error) {
	participants, err := conversation.NewDirectConversationParticipantsPair(participantOneID, participantTwoID)
	if err != nil {
		return nil, err
	}

	cnv, err := d.conversationRepository.GetDirectBetween(ctx, participants)
	if err == nil {
		return cnv, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	cnv = conversation.New(conversation.ConversationTypeDirect)

	err = d.conversationRepository.SaveDirect(ctx, cnv, participants)
	if err != nil {
		return nil, err
	}

	for _, participantID := range participants.List() {
		exists, err := d.participantLookup.Exists(ctx, participantID) 
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, apperr.NewNotFound("Participant doesn't exist")
		}

		membership := membership.New(
			cnv.ID(),
			participantID,
			membership.RoleMember,
		)
		err = d.membershipRepository.Save(ctx, membership)
		if err != nil {
			return nil, err
		}
	}

	return cnv, nil
}
