package command

import (
	"context"
)

type Command interface {
	IsCommand()
	Name() string
}

type Handler[C Command, R any] interface {
	Handle(context.Context, C) (R, error)
}

type HandlerFunc[C Command, R any] func(context.Context, C) (R, error)

func (f HandlerFunc[C, R]) Handle(ctx context.Context, cmd C) (R, error) {
	return f(ctx, cmd)
}
