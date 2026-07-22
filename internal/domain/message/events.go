package message

type MessageSent struct {
	MessageID string
	ChatID    string
}

func (m MessageSent) IsEvent() {}

func (m MessageSent) Name() string {
	return "MessageSent"
}
