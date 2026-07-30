package commandbus

import (
	"context"
	"fmt"
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

func (b *Bus) Dispatch(ctx context.Context, command Command) (any, error) {
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

func RegisterHandler[C Command](b *Bus, handler ApCommandHandler[C]) error {
	c := *new(C)

	var adapted CommandHandler = ApHandlerAdapter[C]{handler: handler}

	for _, m := range b.middlewares {
		adapted = m(adapted)
	}
	if _, ok := b.handlers[c.Name()]; ok {
		return fmt.Errorf("для команды '%s' уже зарегистрирован обработчик", c.Name())
	}
	b.handlers[c.Name()] = adapted
	return nil
}

type Command interface {
	IsCommand()
	Name() string
}

type CommandHandler interface {
	Handle(context.Context, Command) (any, error)
}

type ApCommandHandler[C Command] interface {
	Handle(context.Context, C) (any, error)
}

type ApHandlerAdapter[C Command] struct {
	handler ApCommandHandler[C]
}

func (a ApHandlerAdapter[C]) Handle(ctx context.Context, c Command) (any, error) {
	typedCommand, ok := c.(C)
	if !ok {
		return nil, fmt.Errorf("command handler expected command of type %T, got %T", new(C), c)
	}

	return a.handler.Handle(ctx, typedCommand)
}

type HandlerFunc func(ctx context.Context, command Command) (any, error)

func (f HandlerFunc) Handle(ctx context.Context, command Command) (any, error) {
	return f(ctx, command)
}

type Middleware func(next CommandHandler) CommandHandler
