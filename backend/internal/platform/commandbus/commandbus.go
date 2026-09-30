package commandbus

import (
	"context"
	"fmt"
	"messenger/messenger/internal/shared/application/command"
)

type Bus struct {
	handlers    map[string]CommandHandler
	middlewares []Middleware
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string]CommandHandler),
	}
}

func (b *Bus) Dispatch[R any](ctx context.Context, c command.Command) (R, error) {
	var zero R

	handler, ok := b.handlers[c.Name()]
	if !ok {
		return zero, fmt.Errorf("command bus: not found handler for command '%s'", c.Name())
	}

	response, err := handler.Handle(ctx, c)
	if err != nil {
		return zero, err
	}

	typedResponse, ok := response.(R)
	if !ok {
		return zero, fmt.Errorf("command bus: handler of dispatched command responded with unexpected type: expected %T, got %T", new(R), response)
	}

	return typedResponse, nil
}

func (b *Bus) RegisterHandler[C command.Command, R any](handler command.Handler[C, R]) error {
	c := *new(C)

	var adapted CommandHandler = &AppCommandHandlerAdapter[C, R]{handler: handler}

	if _, ok := b.handlers[c.Name()]; ok {
		return fmt.Errorf("для команды '%s' уже зарегистрирован обработчик", c.Name())
	}
	for _, m := range b.middlewares {
		adapted = m(adapted)
	}
	b.handlers[c.Name()] = adapted
	return nil
}

func (b *Bus) Use(middleware Middleware) {
	for c, h := range b.handlers {
		b.handlers[c] = middleware(h)
	}

	b.middlewares = append(b.middlewares, middleware)
}

type CommandHandler command.Handler[command.Command, any]

var _ CommandHandler = (*AppCommandHandlerAdapter[command.Command, any])(nil)

type AppCommandHandlerAdapter[C command.Command, R any] struct {
	handler command.Handler[C, R]
}

func (a *AppCommandHandlerAdapter[C, R]) Handle(ctx context.Context, c command.Command) (any, error) {
	typedCommand, ok := c.(C)
	if !ok {
		return nil, fmt.Errorf("command handler expected command of type %T, got %T", new(C), c)
	}

	return a.handler.Handle(ctx, typedCommand)
}

var _ CommandHandler = HandlerFunc(nil)

type HandlerFunc func(ctx context.Context, command command.Command) (any, error)

func (f HandlerFunc) Handle(ctx context.Context, command command.Command) (any, error) {
	return f(ctx, command)
}

type Middleware func(next CommandHandler) CommandHandler
