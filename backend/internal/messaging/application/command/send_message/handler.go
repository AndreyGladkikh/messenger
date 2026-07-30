package send_message

import (
	"context"
	"messenger/messenger/internal/messaging/domain/chat"
	"messenger/messenger/internal/messaging/domain/message"
)

type Handler struct {
	messageRepository message.Repository
	chatRepository    chat.Repository
}

func NewHandler(
	messageRepository message.Repository,
	chatRepository chat.Repository,
) *Handler {
	return &Handler{
		messageRepository: messageRepository,
		chatRepository:    chatRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, command *Command) (response any, err error) {
	c, err := h.chatRepository.Get(ctx, command.ChatID)
	if err != nil {
		return nil, err
	}
	if c.IsDeleted() {
		return nil, chat.ErrDeleted
	}

	message := message.Send(
		command.ChatID,
		command.SenderID,
		command.MessageBody,
		command.ReplyToMessageID,
	)

	err = h.messageRepository.Add(ctx, message)

	return nil, err
}
