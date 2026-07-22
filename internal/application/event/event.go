package event

import (
	"context"
	"messenger/messenger/internal/domain"
)

type Handler interface {
	Handle(context.Context, domain.Event) error
	Name() string
}