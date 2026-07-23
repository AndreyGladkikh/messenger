package chat

type ChatCreated struct {
	ChatID string
}

func (c ChatCreated) IsEvent() {}

func (c ChatCreated) Name() string {
	return "ChatCreated"
}