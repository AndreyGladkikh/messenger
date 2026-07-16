package event

import "context"

type Bus struct {
	handlers    map[string][]Handler
	middlewares []Middleware
}

func NewBus() *Bus {
	return &Bus{}
}

func (b *Bus) Register(ctx context.Context, e Event, h Handler) {
	for _, m := range b.middlewares {
		h = m(h)
	}

	eHandlers, ok := b.handlers[e.Name()]
	if ok {
		eHandlers = append(eHandlers, h)
	} else {
		eHandlers = []Handler{h}
	}
	b.handlers[e.Name()] = eHandlers
}

func (b *Bus) Dispatch(ctx context.Context, e Event) error {
	if handlers, ok := b.handlers[e.Name()]; ok {
		for _, h := range handlers {
			if err := h(ctx, e); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (b *Bus) Use(m Middleware) {
	b.middlewares = append(b.middlewares, m)

	for eventName, eventHandlers := range b.handlers {
		for idx, h := range eventHandlers {
			b.handlers[eventName][idx] = m(h)
		}
	}
}

type Middleware func(Handler) Handler

type Event interface {
	Name() string
}

// type Handler interface {
// 	Handle(ctx context.Context, e any) error
// }

type Handler func(ctx context.Context, e Event) error
