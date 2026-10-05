package service

import (
	"context"
	"fmt"
	"messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"messenger/messenger/internal/conversations/domain/model/membership"
	"messenger/messenger/internal/conversations/domain/model/message"
	"uuid"
)

var ErrSenderNotParticipantOfConversation = fmt.Errorf("%w: sender is not a participant in conversation", domain.ErrForbidden)

type MessageSender struct {
	membershipRepository membership.Repository
	messagesRepository     message.Repository
}

func NewMessageSender(
	membershipRepository membership.Repository,
	messagesRepository message.Repository,
) *MessageSender {
	return &MessageSender{
		membershipRepository: membershipRepository,
		messagesRepository:     messagesRepository,
	}
}

func (s *MessageSender) Send(
	ctx context.Context,
	sender uuid.UUID,
	cnv *conversation.Conversation,
	body string,
) (*message.Message, error) {
	err := s.ensureSenderHasAccessToConversation(ctx, sender, cnv)
	if err != nil {
		return nil, err
	}

	msg := message.New(
		sender,
		cnv.ID(),
		body,
	)

	err = s.messagesRepository.Save(ctx, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageSender) ensureSenderHasAccessToConversation(ctx context.Context, senderID uuid.UUID, cnv *conversation.Conversation) error {
	exists, err := s.membershipRepository.Exists(ctx, cnv.ID(), senderID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrSenderNotParticipantOfConversation
	}

	// switch cnv.Type {
	// case conversation.ConversationTypeDirect:
	// 	exists, err := s.conversationRepository.ExistsDirectWith(ctx, cnv.ID(), senderID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !exists {
	// 		return ErrSenderNotParticipantOfConversation
	// 	}
	// }
	return nil
}
