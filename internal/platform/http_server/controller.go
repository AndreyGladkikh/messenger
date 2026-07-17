package http_server

import (
	"encoding/json"
	"messenger/messenger/internal/application/command/send_message"
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

// func (c *Controller) successResponse(w http.ResponseWriter, _ *http.Request, v any, httpStatus int) {
// 	response := map[string]any{
// 		"status": "success",
// 		"data":   v,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(httpStatus)

// 	e := json.NewEncoder(w).Encode(response)
// 	if e != nil {
// 		http.Error(w, e.Error(), http.StatusInternalServerError)
// 	}
// }

// func (c *Controller) errorResponse(w http.ResponseWriter, _ *http.Request, err error) {
// 	response := map[string]any{
// 		"status": "error",
// 		"error":  err,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(getErrorStatus(err))

// 	e := json.NewEncoder(w).Encode(response)
// 	if e != nil {
// 		http.Error(w, e.Error(), http.StatusInternalServerError)
// 	}
// }

func (c *Controller) sendMessage(w http.ResponseWriter, r *http.Request) {
	var request Message
	json.NewDecoder(r.Body).Decode(&request)

	command := &send_message.Command{
		ChatID:           request.ChatID,
		MessageBody:      request.Body,
		ReplyToMessageID: request.ReplyToMessageID,
		Attachments:      request.Attachments,
	}
	response, err := c.commandBus.Dispatch(r.Context(), command)

	NewResponse(response, err).WriteTo(w)
}

// func (c *Controller) sendMessage(w http.ResponseWriter, r *http.Request) {
// 	var request Message
// 	json.NewDecoder(r.Body).Decode(&request)

// 	command := &send_message.Command{
// 		ChatID:           request.ChatID,
// 		MessageBody:      request.Body,
// 		ReplyToMessageID: request.ReplyToMessageID,
// 		Attachments:      request.Attachments,
// 	}
// 	response, err := c.commandBus.Dispatch(r.Context(), command)

// 	if err != nil {
// 		c.errorResponse(w, r, err)
// 	}
// 	c.successResponse(w, r, response, http.StatusOK)
// }
