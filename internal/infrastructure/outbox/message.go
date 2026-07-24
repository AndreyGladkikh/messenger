package outbox

import "messenger/messenger/internal/domain"

type Message struct {
	Event *Event
}

type Event struct {
	ID string
	domain.Event
}