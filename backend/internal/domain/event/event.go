package event

type DomainEvent interface {
	Name() string
	// isEvent()
	// isDomainEvent()
}
