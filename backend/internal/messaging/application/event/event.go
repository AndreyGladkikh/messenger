package event

import (
	"context"
	"messenger/messenger/internal/messaging/domain"
)

type Handler[E domain.Event] interface {
	Handle(context.Context, E) error
	Name() string
}
