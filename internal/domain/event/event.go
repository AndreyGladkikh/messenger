package event

type DomainEvent interface {
	Name() string
	// isDomainEvent()
}