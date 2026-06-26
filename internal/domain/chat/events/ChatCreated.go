package events

import "time"

type ChatCreated struct {
	chatID string
	occuredAt *time.Time
}

func NewChatCreated(
	chatID string,
	occuredAt *time.Time,
) ChatCreated {
	return ChatCreated{
		chatID: chatID,
		occuredAt: occuredAt,
	}
}

func (c ChatCreated) isDomainEvent() {}