package domain

type Event interface {
	IsEvent()
	Name() string
}
