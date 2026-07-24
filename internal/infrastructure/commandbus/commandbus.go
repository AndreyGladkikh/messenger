package commandbus

import (
	"context"
	"fmt"
	"messenger/messenger/internal/application/command"
)

type Bus struct {
	handlers    map[string]Handler
	middlewares []Middleware
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string]Handler),
	}
}

func (b *Bus) Register(command command.Command, handler Handler) error {
	for _, m := range b.middlewares {
		handler = m(handler)
	}
	if _, ok := b.handlers[command.Name()]; ok {
		return fmt.Errorf("для команды '%s' уже зарегистрирован обработчик", command.Name())
	}
	b.handlers[command.Name()] = handler
	return nil
}

func (b *Bus) Dispatch(ctx context.Context, command command.Command) (any, error) {
	handler, ok := b.handlers[command.Name()]
	if !ok {
		var zero any
		return zero, fmt.Errorf("не найден обработчик команды '%s'", command.Name())
	}

	return handler.Handle(ctx, command)
}

func (b *Bus) Use(middleware Middleware) {
	for c, h := range b.handlers {
		b.handlers[c] = middleware(h)
	}

	b.middlewares = append(b.middlewares, middleware)
}

func RegisterHandler[C command.Command](b *Bus, command C, handler command.Handler[C]) error {
	var adapted Handler = CommandHandlerAdapter[C]{handler: handler}

	for _, m := range b.middlewares {
		adapted = m(adapted)
	}
	if _, ok := b.handlers[command.Name()]; ok {
		return fmt.Errorf("для команды '%s' уже зарегистрирован обработчик", command.Name())
	}
	b.handlers[command.Name()] = adapted
	return nil
}

type Handler interface {
	Handle(context.Context, command.Command) (any, error)
}

type CommandHandlerAdapter[C command.Command] struct {
	handler command.Handler[C]
}

func (ha CommandHandlerAdapter[C]) Handle(ctx context.Context, c command.Command) (any, error) {
	typedCommand, ok := c.(C)
	if !ok {
		return nil, fmt.Errorf("command handler expected another command type")
	}

	return ha.handler.Handle(ctx, typedCommand)
}

type HandlerFunc func(ctx context.Context, command command.Command) (any, error)

func (f HandlerFunc) Handle(ctx context.Context, command command.Command) (any, error) {
	return f(ctx, command)
}

type Middleware func(next Handler) Handler