package pubsub

import (
	"context"
	"fmt"

	"uuid"
)

type PersonalEventsPubSub struct {
	pubSub *PubSub
}

func NewPersonalEventsPubSub(pubSub *PubSub) *PersonalEventsPubSub {
	return &PersonalEventsPubSub{
		pubSub: pubSub,
	}
}

func (u *PersonalEventsPubSub) Publish(ctx context.Context, userID uuid.UUID, msg Message) error {
	return u.pubSub.Publish(ctx, ChannelNamePersonalEvents(userID), msg)
}

func (u *PersonalEventsPubSub) Subscription(ctx context.Context, subscriberName string, userID uuid.UUID) *Subscription {
	return u.pubSub.Subscribe(ctx, subscriberName, ChannelNamePersonalEvents(userID))
}

func ChannelNamePersonalEvents(userID uuid.UUID) string {
	return fmt.Sprintf("users:%s:events", userID.String())
}
