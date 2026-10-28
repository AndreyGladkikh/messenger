package pubsub

import (
	"context"
	"fmt"

	"uuid"
)

type UserEventsPubSub struct {
	PubSub *PubSub
}

func (u *UserEventsPubSub) Publish(ctx context.Context, userID uuid.UUID, msg Message) error {
	return u.PubSub.Publish(ctx, ChannelNameUserEvents(userID), msg)
}

func (u *UserEventsPubSub) Subscription(ctx context.Context, subscriberName string, userID uuid.UUID) *Subscription {
	return u.PubSub.Subscribe(ctx, subscriberName, ChannelNameUserEvents(userID))
}

func ChannelNameUserEvents(userID uuid.UUID) string {
	return fmt.Sprintf("users:%s:events", userID.String())
}
