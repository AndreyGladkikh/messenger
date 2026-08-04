package domain

import "github.com/google/uuid"

type Aggregate interface {
	ID() uuid.UUID
	AddEvent(Event)
	PullEvents() []Event
}

type BaseAggregate struct {
	events []Event
}

func (a *BaseAggregate) AddEvent(e Event) {
	a.events = append(a.events, e)
}

func (a *BaseAggregate) PullEvents() []Event {
	events := a.events
	a.events = nil
	return events
}
