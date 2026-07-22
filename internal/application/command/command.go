package command

import (
	"context"
	"errors"
)

var ErrValidation = errors.New("validation error")

type Command interface {
	IsCommand()
	Name() string
}

type Handler[C Command] interface {
	Handle(context.Context, C) (any, error)
}