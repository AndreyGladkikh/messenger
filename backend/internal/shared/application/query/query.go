package query

import (
	"context"
	"errors"
)

var ErrValidation = errors.New("validation error")

type Query interface {
	IsQuery()
	Name() string
}

type Handler[Q Query, R any] interface {
	Handle(context.Context, Q) (R, error)
}
