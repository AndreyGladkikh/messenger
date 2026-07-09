package message

type MessageCreated struct {
	MessageID string
}

func (m MessageCreated) Name() string {
	return "event.message.created"
}

type MessageSent struct {
	MessageID string
	ChatID string
}

func (m MessageSent) Name() string {
	return "event.message.sent"
}