package query

import (
	"context"
)

type Query interface {
	Name() string
}

type Handler[Q Query, R any] interface {
	Handle(context.Context, Q) (R, error)
}

type HandlerFunc[Q Query, R any] func(context.Context, Q) (R, error)

func (f HandlerFunc[Q, R]) Handle(ctx context.Context, qry Q) (R, error) {
	return f(ctx, qry)
}
