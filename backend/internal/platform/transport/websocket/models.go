package websocket

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"uuid"
)

type WSMessage struct {
	ID      string         `json:"id,omitempty"`
	Type    WSMessageType  `json:"type"`
	Payload jsontext.Value `json:"data,omitempty"`
	Error   jsontext.Value `json:"error,omitempty"`
}

func NewWSMessage(typ WSMessageType, opts ...wsMessageOption) (WSMessage, error) {
	m := WSMessage{
		Type: typ,
	}

	for _, opt := range opts {
		err := opt(&m)
		if err != nil {
			return WSMessage{}, err
		}
	}

	return m, nil
}

func newResponse(id string, data any, err error) (WSMessage, error) {
	opts := []wsMessageOption{withID(id)}
	if data != nil {
		opts = append(opts, withData(data))
	}
	if err != nil {
		opts = append(opts, withError(err))
	}

	return NewWSMessage(
		wsMessageTypeResponse,
		opts...,
	)
}

type WSMessageType string

const (
	wsMessageTypeResponse WSMessageType = "response"

	wsMessageTypeSendMessage            WSMessageType = "sendMessage"
	wsMessageTypeSendFirstDirectMessage WSMessageType = "sendFirstDirectMessage"
	wsMessageTypeAddSubscription        WSMessageType = "addSubscription"

	wsMessageTypeMessageSent WSMessageType = "messageSent"
)

type wsMessageOption func(*WSMessage) error

func withID(id string) wsMessageOption {
	return func(m *WSMessage) error {
		m.ID = id
		return nil
	}
}

func withData(data any) wsMessageOption {
	return func(m *WSMessage) error {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		m.Payload = jsonData
		return nil
	}
}

type SendFirstDirectMessagePayload struct {
	Body        string    `json:"body"`
	RecipientID uuid.UUID `json:"recipient"`
}
