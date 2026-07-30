package http_server

import (
	"encoding/json"
	"messenger/messenger/internal/auth/application/command/register"
	"messenger/messenger/internal/messaging/application/command/create_private_chat"
	"messenger/messenger/internal/messaging/application/command/send_message"
	"messenger/messenger/internal/messaging/infrastructure/auth"
	"messenger/messenger/internal/platform/commandbus"
	"net/http"
)

type Controller struct {
	commandBus *commandbus.Bus
}

func NewController(
	commandBus *commandbus.Bus,
) *Controller {
	return &Controller{
		commandBus: commandBus,
	}
}

func (c *Controller) registerUser(w http.ResponseWriter, r *http.Request) {
	var request RegisterUserRequest
	json.NewDecoder(r.Body).Decode(&request)

	command := &register.Command{
		Login:    request.Login,
		Password: request.Password,
	}
	response, err := c.commandBus.Dispatch(r.Context(), command)

	NewResponse(response, err).WriteTo(w)
}

func (c *Controller) sendMessage(w http.ResponseWriter, r *http.Request) {
	var request Message
	json.NewDecoder(r.Body).Decode(&request)

	senderID, _ := auth.UserIDFromContext(r.Context())

	command := &send_message.Command{
		SenderID:         senderID,
		ChatID:           request.ChatID,
		MessageBody:      request.Body,
		ReplyToMessageID: request.ReplyToMessageID,
		Attachments:      request.Attachments,
	}
	response, err := c.commandBus.Dispatch(r.Context(), command)

	NewResponse(response, err).WriteTo(w)
}

func (c *Controller) createPrivateChat(w http.ResponseWriter, r *http.Request) {
	var request CreatePrivateChatRequest
	json.NewDecoder(r.Body).Decode(&request)

	userID, _ := auth.UserIDFromContext(r.Context())

	command := &create_private_chat.Command{
		InitiatorID:    userID,
		ChatWithUserID: request.ChatWithUserID,
	}
	response, err := c.commandBus.Dispatch(r.Context(), command)

	NewResponse(response, err).WriteTo(w)
}
