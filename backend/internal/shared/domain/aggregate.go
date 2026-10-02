package domain

import "uuid"

type Aggregate interface {
	ID() uuid.UUID
	AddEvent(Event)
	PullEvents() []Event
}

type AggregateRoot struct {
	events []Event
}

func (a *AggregateRoot) AddEvent(e Event) {
	a.events = append(a.events, e)
}

func (a *AggregateRoot) PullEvents() []Event {
	events := a.events
	a.events = nil
	return events
}
