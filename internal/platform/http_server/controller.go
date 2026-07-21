package http_server

import (
	"encoding/json"
	"messenger/messenger/internal/application/command/send_message"
	"messenger/messenger/internal/platform/auth"
	"messenger/messenger/internal/platform/command_bus"
	"net/http"
)

type Controller struct {
	commandBus *command_bus.Bus
}

func NewController(
	commandBus *command_bus.Bus,
) *Controller {
	return &Controller{
		commandBus: commandBus,
	}
}

func (c *Controller) sendMessage(w http.ResponseWriter, r *http.Request) {
	var request Message
	json.NewDecoder(r.Body).Decode(&request)

	senderID, _ := auth.UserIDFromContext(r.Context())

	command := &send_message.Command{
		SenderID: senderID,
		ChatID:           request.ChatID,
		MessageBody:      request.Body,
		ReplyToMessageID: request.ReplyToMessageID,
		Attachments:      request.Attachments,
	}
	response, err := c.commandBus.Dispatch(r.Context(), command)

	NewResponse(response, err).WriteTo(w)
}

// func (c *Controller) createChat(w http.ResponseWriter, r *http.Request) {
// 	var request Chat
// 	json.NewDecoder(r.Body).Decode(&request)

// 	command := &create_chat.Command{
// 		Type: request.Type,
// 		ChatName: request.Name,
// 	}
// 	response, err := c.commandBus.Dispatch(r.Context(), command)

// 	NewResponse(response, err).WriteTo(w)
// }