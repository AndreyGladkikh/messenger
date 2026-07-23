package domain

type Aggregate interface {
	ID() string
	// AggregateName() string
	// AggregateID() string
	AddEvent(Event)
	PullEvents() []Event
}

type BaseAggregate struct {
	events []Event
}

// func (a *BaseAggregate) AggregateID() string {
// 	return fmt.Sprintf("%s:%s", a.AggregateName(), a.ID())
// }

func (a *BaseAggregate) AddEvent(e Event) {
	a.events = append(a.events, e)
}

func (a *BaseAggregate) PullEvents() []Event {
	events := a.events
	a.events = nil
	return events
}
