package pubsub

import (
	"context"
	"fmt"

	"uuid"
)

type UserEventsPubSub struct {
	pubSub *PubSub
}

func NewUserEventsPubSub(pubSub *PubSub) *UserEventsPubSub {
	return &UserEventsPubSub{
		pubSub: pubSub,
	}
}

func (u *UserEventsPubSub) Publish(ctx context.Context, userID uuid.UUID, msg Message) error {
	return u.pubSub.Publish(ctx, ChannelNameUserEvents(userID), msg)
}

func (u *UserEventsPubSub) Subscription(ctx context.Context, subscriberName string, userID uuid.UUID) *Subscription {
	return u.pubSub.Subscribe(ctx, subscriberName, ChannelNameUserEvents(userID))
}

func ChannelNameUserEvents(userID uuid.UUID) string {
	return fmt.Sprintf("users:%s:events", userID.String())
}
