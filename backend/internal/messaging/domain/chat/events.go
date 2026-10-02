package chat

import "uuid"

const ChatCreatedEventName = "ChatCreated"

type ChatCreated struct {
	ChatID uuid.UUID
}

func (c ChatCreated) IsEvent() {}

func (c ChatCreated) Name() string {
	return ChatCreatedEventName
}
