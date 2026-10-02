package service

import (
	"context"
	"errors"
	"messenger/messenger/internal/conversation/domain/model/conversation"
	"messenger/messenger/internal/conversation/domain/model/membership"
	"messenger/messenger/internal/conversation/domain/model/participant"
	"messenger/messenger/internal/shared/domain"
)

type DirectConversationResolver struct {
	conversationRepository conversation.Repository
	membershipRepository   membership.Repository
}

func (f *DirectConversationResolver) FindOrCreate(
	ctx context.Context,
	participantOne participant.Participant,
	participantTwo participant.Participant,
) (*conversation.Conversation, error) {
	participants, err := conversation.NewDirectConversationParticipantsPair(participantOne, participantTwo)
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

	err = f.conversationRepository.Save(ctx, cnv)
	if err != nil {
		return nil, err
	}

	for _, participant := range participants.List() {
		membership := membership.New(
			cnv.ID(),
			participant,
			membership.RoleMember,
		)
		err = f.membershipRepository.Save(ctx, membership)
		if err != nil {
			return nil, err
		}
	}

	return cnv, nil
}
