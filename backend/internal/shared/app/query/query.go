package query

import (
	"context"
)

type Query interface {
	IsQuery()
	Name() string
}

type Handler[Q Query, R any] interface {
	Handle(context.Context, Q) (R, error)
}
