package message

type MessageCreated struct {
	MessageID string
}

func (m MessageCreated) Name() string {
	return "event.message.created"
}