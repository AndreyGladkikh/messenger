package events

import "time"

type ChatCreated struct {
	chatID string
	userID string
	occuredAt *time.Time
}

func NewChatCreated(
	chatID string,
	userID string,
	occuredAt *time.Time,
) ChatCreated {
	return ChatCreated{
		chatID: chatID,
		userID: userID,
		occuredAt: occuredAt,
	}
}

func (c ChatCreated) isDomainEvent() {}