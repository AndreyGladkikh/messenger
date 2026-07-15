package command

import (
	"context"
	"fmt"
)

type Bus struct {
	handlers    map[string]Handler
	middlewares []HandlerMiddleware
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string]Handler),
	}
}

func (b *Bus) Register(command Command, handler Handler) error {
	for _, m := range b.middlewares {
		handler = m(handler)
	}
	if _, ok := b.handlers[command.Name()]; ok {
		return fmt.Errorf("для команды '%s' уже зарегистрирован обработчик", command.Name())
	}
	b.handlers[command.Name()] = handler
	return nil
}

func (b *Bus) Dispatch(ctx context.Context, command Command) (any, error) {
	handler, ok := b.handlers[command.Name()]
	if !ok {
		var zero any
		return zero, fmt.Errorf("не найден обработчик команды '%s'", command.Name())
	}

	return handler.Handle(ctx, command)
}

func (b *Bus) Use(middleware HandlerMiddleware) {
	b.middlewares = append(b.middlewares, middleware)
}
