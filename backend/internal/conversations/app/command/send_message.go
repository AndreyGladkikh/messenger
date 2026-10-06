package command

import (
	"context"
	"messenger/messenger/internal/conversations/app/service"
	"messenger/messenger/internal/conversations/domain/model/conversation"
	"uuid"
)

type SendMessage struct {
	SenderID       uuid.UUID
	ConversationID uuid.UUID
	Body           string
}

func (cmd SendMessage) IsCommand() {}

func (cmd SendMessage) Name() string {
	return "SendMessage"
}

type SendMessageHandler struct {
	conversationRepository conversation.Repository
	messageSender          *service.MessageSender
}

func NewSendMessageHandler(
	conversationRepository conversation.Repository,
	messageSender *service.MessageSender,
) *SendMessageHandler {
	return &SendMessageHandler{
		conversationRepository: conversationRepository,
		messageSender:          messageSender,
	}
}

func (h *SendMessageHandler) Handle(ctx context.Context, cmd SendMessage) (any, error) {
	cnv, err := h.conversationRepository.Get(ctx, cmd.ConversationID)
	if err != nil {
		return nil, err
	}

	_, err = h.messageSender.Send(ctx, cmd.SenderID, cnv, cmd.Body)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
