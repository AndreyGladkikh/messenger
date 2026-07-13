package event

type Bus struct {
}

func NewBus() *Bus {
	return &Bus{}
}

func (b *Bus) Dispatch(e any) error {
	return nil
}
