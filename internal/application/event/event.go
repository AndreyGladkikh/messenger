package event

import (
	"context"
	"messenger/messenger/internal/domain"
)

type Handler[T domain.Event] interface {
	Handle(context.Context, T) error
	Name() string
}