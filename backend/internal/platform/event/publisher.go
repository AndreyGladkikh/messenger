package event

import (
	"context"
	"messenger/messenger/internal/platform/uow"
	"messenger/messenger/internal/shared/app/event"
)

var _ event.Publisher = (*Publisher)(nil)

type Publisher struct {
}

func (p *Publisher) Publish(ctx context.Context, event event.Event) {
	uw, ok := uow.FromContext(ctx)
	if !ok {
		return
	}

	uw.AddEvent(event)
}
