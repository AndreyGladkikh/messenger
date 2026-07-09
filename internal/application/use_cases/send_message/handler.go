package send_message

import (
	"context"
	"messenger/messenger/internal/application/user_id"
	"messenger/messenger/internal/domain/event"
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

	userID, _ := user_id.FromContext(ctx)

	message, messageCreated := message.NewMessage(
		h.messageRepository.NextID(),
		command.MessageBody,
		command.ChatID,
		userID,
		command.RecipientID,
		command.ReplyToMessageID,
	)

	event.Publisher().Publish(messageCreated)

	err = h.messageRepository.Add(ctx, message)

	return nil, err
}
