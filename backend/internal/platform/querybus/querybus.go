package querybus

import (
	"context"
	"fmt"
	query1 "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/shared/app/query"
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

func (b *Bus) Dispatch[R any](ctx context.Context, q query.Query) (R, error) {
	var zero R

	handler, ok := b.handlers[q.Name()]
	if !ok {
		return zero, fmt.Errorf("query bus: not found handler for query '%s'", q.Name())
	}

	response, err := handler.Handle(ctx, q)
	if err != nil {
		return zero, err
	}

	typedResponse, ok := response.(R)
	if !ok {
		return zero, fmt.Errorf("query bus: handler of dispatched command responded with unexpected type: expected %T, got %T", new(R), response)
	}

	return typedResponse, nil
}

func (b *Bus) RegisterHandler[Q query.Query, R any](handler query.Handler[Q, R]) error {
	var q Q

	var adapted QueryHandler = AppQueryHandlerAdapter[Q, R]{handler: handler}

	if _, ok := b.handlers[q.Name()]; ok {
		return fmt.Errorf("query bus: already has handler for query '%s'", q.Name())
	}
	for _, m := range b.middlewares {
		adapted = m(adapted)
	}
	b.handlers[q.Name()] = adapted
	return nil
}

func (b *Bus) Use(middleware Middleware) {
	for c, h := range b.handlers {
		b.handlers[c] = middleware(h)
	}

	b.middlewares = append(b.middlewares, middleware)
}

type QueryHandler query.Handler[query.Query, any]

var _ QueryHandler = AppQueryHandlerAdapter[query.Query, any]{}

type AppQueryHandlerAdapter[Q query.Query, R any] struct {
	handler query.Handler[Q, R]
}

func (a AppQueryHandlerAdapter[Q, R]) Handle(ctx context.Context, q query.Query) (any, error) {
	var qTP Q

	typedQuery, ok := q.(Q)
	if !ok {
		return nil, fmt.Errorf("query handler expected query of type %T, got %T", qTP, q)
	}

	return a.handler.Handle(ctx, typedQuery)
}

var _ QueryHandler = HandlerFunc(nil)

type HandlerFunc func(ctx context.Context, q query.Query) (any, error)

func (f HandlerFunc) Handle(ctx context.Context, q query.Query) (any, error) {
	return f(ctx, q)
}

type Middleware func(next QueryHandler) QueryHandler

func InitBus(
	getAccountHandler *query1.GetAccountHandler,
) *Bus {
	b := NewBus()

	b.RegisterHandler(getAccountHandler)

	return b
}
