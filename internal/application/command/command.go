package command

import "context"

type Command interface {
	IsCommand()
	Name() string
}

type Handler interface {
	Handle(ctx context.Context, command Command) (any, error)
}

type HandlerMiddleware func(next Handler) Handler

type HandlerFunc func(ctx context.Context, command Command) (any, error)

func (f HandlerFunc) Handle(ctx context.Context, command Command) (any, error) {
	return f(ctx, command)
}

// type HandlerFuncMiddleware func(next HandlerFunc) HandlerFunc
