package command_bus

import (
	"context"
	"fmt"
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

func (b *Bus) Use(middleware Middleware) {
	for c, h := range b.handlers {
		b.handlers[c] = middleware(h)
	}

	b.middlewares = append(b.middlewares, middleware)
}

type Command interface {
	IsCommand()
	Name() string
}

type Handler interface {
	Handle(context.Context, Command) (any, error)
}

type HandlerFunc func(ctx context.Context, command Command) (any, error)

func (f HandlerFunc) Handle(ctx context.Context, command Command) (any, error) {
	return f(ctx, command)
}

type Middleware func(next Handler) Handler


// import (
// 	"context"
// 	"messenger/messenger/internal/adapters/out/postgres/transaction"
// 	"messenger/messenger/internal/application/command/send_message"
// 	"messenger/messenger/internal/platform/command"
// 	"messenger/messenger/internal/platform/command/bus_middlewares"
// )

// func BuildCommandBus(
// 	txManager *transaction.Manager,
// 	sendMessageHandler *send_message.Handler,
// ) *command.Bus {
// 	bus := new(command.Bus)

// 	txMiddlewareContainer := bus_middlewares.NewTransactionMiddlewareContainer(txManager)

// 	bus.Use(txMiddlewareContainer.Middleware)

// 	bus.Register(&send_message.Command{}, command.HandlerFunc(func(ctx context.Context, cmd command.Command) (any, error) {
// 		return sendMessageHandler.Handle(ctx, cmd.(*send_message.Command))
// 	}))

// 	return bus
// }
