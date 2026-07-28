package chat

const ChatCreatedEventName = "ChatCreated"

type ChatCreated struct {
	ChatID string
}

func (c ChatCreated) IsEvent() {}

func (c ChatCreated) Name() string {
	return ChatCreatedEventName
}
