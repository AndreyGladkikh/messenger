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
