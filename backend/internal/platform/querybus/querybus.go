package querybus

import (
	"context"
	"fmt"
	"messenger/messenger/internal/messaging/application/query/get_chat_list"
	"messenger/messenger/internal/shared/application/query"
)

type Bus struct {
	handlers    map[string]QueryHandler
	middlewares []Middleware
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string]QueryHandler),
	}
}

func (b *Bus) Dispatch(ctx context.Context, q query.Query) (any, error) {
	handler, ok := b.handlers[q.Name()]
	if !ok {
		var zero any
		return zero, fmt.Errorf("не найден обработчик запроса '%s'", q.Name())
	}

	return handler.Handle(ctx, q)
}

func (b *Bus) Use(middleware Middleware) {
	for c, h := range b.handlers {
		b.handlers[c] = middleware(h)
	}

	b.middlewares = append(b.middlewares, middleware)
}

func RegisterHandler[Q query.Query, R any](b *Bus, handler query.Handler[Q, R]) error {
	q := *new(Q)

	var adapted QueryHandler = ApQueryHandlerAdapter[Q, R]{handler: handler}

	for _, m := range b.middlewares {
		adapted = m(adapted)
	}
	if _, ok := b.handlers[q.Name()]; ok {
		return fmt.Errorf("для запроса '%s' уже зарегистрирован обработчик", q.Name())
	}
	b.handlers[q.Name()] = adapted
	return nil
}

type QueryHandler interface {
	Handle(context.Context, query.Query) (any, error)
}

type ApQueryHandlerAdapter[Q query.Query, R any] struct {
	handler query.Handler[Q, R]
}

func (a ApQueryHandlerAdapter[Q, R]) Handle(ctx context.Context, q query.Query) (any, error) {
	typedQuery, ok := q.(Q)
	if !ok {
		return nil, fmt.Errorf("query handler expected query of type %T, got %T", new(Q), q)
	}

	return a.handler.Handle(ctx, typedQuery)
}

type HandlerFunc func(ctx context.Context, q query.Query) (any, error)

func (f HandlerFunc) Handle(ctx context.Context, q query.Query) (any, error) {
	return f(ctx, q)
}

type Middleware func(next QueryHandler) QueryHandler

func InitBus(
	getChatListHandler *get_chat_list.Handler,
) *Bus {
	b := NewBus()

	RegisterHandler(b, getChatListHandler)

	return b
}