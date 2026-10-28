package pubsub

import (
	"context"
	"fmt"

	"uuid"
)

// todo rename to ConversationEvents
type ChatEventsPubSub struct {
	pubSub *PubSub
}

func NewChatEventsPubSub(
	pubSub *PubSub,
) *ChatEventsPubSub {
	return &ChatEventsPubSub{
		pubSub: pubSub,
	}
}

func (c *ChatEventsPubSub) Publish(ctx context.Context, chatID uuid.UUID, msg Message) error {
	return c.pubSub.Publish(ctx, ChannelNameConversationEvents(chatID), msg)
}

func (c *ChatEventsPubSub) Subscription(ctx context.Context, subscriberName string) (*ChatEventsSubscription, error) {
	subscription := c.pubSub.Subscription(subscriberName)

	return &ChatEventsSubscription{subscription}, nil
}

type ChatEventsSubscription struct {
	*Subscription
}

func (s *ChatEventsSubscription) AddChats(ctx context.Context, chatIDs ...uuid.UUID) {
	channels := make([]string, 0, len(chatIDs))
	for _, cid := range chatIDs {
		channels = append(channels, ChannelNameConversationEvents(cid))
	}
	s.Subscription.AddChannels(ctx, channels...)
}

func ChannelNameConversationEvents(chatID uuid.UUID) string {
	return fmt.Sprintf("chats:%s:events", chatID.String())
}
