package event

import (
	"context"
	sharedDomain "messenger/messenger/internal/shared/domain"
)

type Handler[E sharedDomain.Event] interface {
	Handle(context.Context, E) error
	Name() string
}
