package message

const MessageSentName = "MessageSent"

type MessageSent struct {
	MessageID string
	ChatID    string
}

func (m MessageSent) IsEvent() {}

func (m MessageSent) Name() string {
	return MessageSentName
}
