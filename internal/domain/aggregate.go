package domain

import (
	// "fmt"
	"messenger/messenger/internal/domain/event"
)

type Aggregate interface {
	ID() string
	// AggregateName() string
	// AggregateID() string
	AddEvent(event.DomainEvent)
	PullEvents() []event.DomainEvent
}

type BaseAggregate struct {
	events []event.DomainEvent
}

// func (a *BaseAggregate) AggregateID() string {
// 	return fmt.Sprintf("%s:%s", a.AggregateName(), a.ID())
// }

func (a *BaseAggregate) AddEvent(e event.DomainEvent) {
	a.events = append(a.events, e)
}

func (a *BaseAggregate) PullEvents() []event.DomainEvent {
	events := a.events
	a.events = nil
	return events
}
