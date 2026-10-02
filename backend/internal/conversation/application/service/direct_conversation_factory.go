package service

import (
	"context"
	"errors"
	"messenger/messenger/internal/conversation/domain/model/conversation"
	"messenger/messenger/internal/conversation/domain/model/membership"
	"messenger/messenger/internal/shared/domain"
)

type DirectConversationResolver struct {
	conversationRepository conversation.Repository
	membershipRepository   membership.Repository
}

func (f *DirectConversationResolver) Create(
	ctx context.Context,
	participants conversation.DirectConversationParticipantsPair,
) (*conversation.Conversation, error) {
	// exists, err := f.conversationRepository.ExistsDirectBetween(ctx, participants)
	// if err != nil {
	// 	return nil, err
	// }

	// if exists {
	// 	return nil, domain.ErrAlreadyExists
	// }

	cnv, err := f.conversationRepository.GetDirectBetween(ctx, participants)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	if cnv != nil {
		return cnv, nil
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
