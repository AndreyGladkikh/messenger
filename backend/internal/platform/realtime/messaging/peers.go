package messaging

import (
	"context"
	"messenger/messenger/internal/platform/pubsub"
)

type Peers struct {
	PubSub pubsub.PubSub
}

func (h *Peers) Handle(ctx context.Context, msg []byte) error {
	// h.PubSub.Publish(ctx, )
	return nil
}
