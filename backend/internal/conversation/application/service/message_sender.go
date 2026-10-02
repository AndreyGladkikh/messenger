package service

import (
	"context"
	"fmt"
	"messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/conversation/domain/model/conversation"
	"messenger/messenger/internal/conversation/domain/model/message"
	"messenger/messenger/internal/conversation/domain/model/participant"
)

var ErrSenderNotParticipantOfConversation = fmt.Errorf("%w: sender is not a participant in conversation", domain.ErrForbidden)

type MessageSender struct {
	conversationRepository conversation.Repository
	messagesRepository     message.Repository
}

func NewMessageSender(
	conversationRepository conversation.Repository,
	messagesRepository message.Repository,
) *MessageSender {
	return &MessageSender{
		conversationRepository: conversationRepository,
		messagesRepository:     messagesRepository,
	}
}

func (s *MessageSender) Send(
	ctx context.Context,
	sender participant.Participant,
	cnv *conversation.Conversation,
	body string,
) (*message.Message, error) {
	err := s.ensureParticipantOfConversation(ctx, sender, cnv)
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

func (s *MessageSender) ensureParticipantOfConversation(ctx context.Context, sender participant.Participant, cnv *conversation.Conversation) error {
	switch cnv.Type {
	case conversation.ConversationTypeDirect:
		exists, err := s.conversationRepository.ExistsDirectWith(ctx, cnv.ID(), sender)
		if err != nil {
			return err
		}
		if !exists {
			return ErrSenderNotParticipantOfConversation
		}
	}
	return nil
}
