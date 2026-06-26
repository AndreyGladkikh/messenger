package send_private_message

import (
	"context"
	"messenger/messenger/internal/domain/message"
)

type Handler struct {
	messageRepository message.Repository
}

func (h *Handler) Handle(ctx context.Context, command *Command) (response any, err error) {
	chatID := command.ChatID
	if chatID == "" {
		chatID = ""
	}

	message := message.NewMessage(
		h.messageRepository.NextID(),
		command.MessageBody,
		// todo getCurrentUserId
		"",
		command.RecipientID,
	)

	err = h.messageRepository.Add(ctx, message)

	return nil, err
}
