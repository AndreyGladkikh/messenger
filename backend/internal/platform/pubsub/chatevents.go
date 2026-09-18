package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"messenger/messenger/internal/shared/application/event"
	"messenger/messenger/internal/shared/domain"

	"github.com/google/uuid"
)

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

func (c *ChatEventsPubSub) Publish(ctx context.Context, chatID uuid.UUID, e event.Envelope[domain.Event]) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to publish message to chat events pubsub: %w", err)
	}

	c.pubSub.Publish(ctx, ChannelNameChatEvents(chatID), payload)

	return nil
}

func (c *ChatEventsPubSub) Subscription(ctx context.Context, subscriberName string, chatIDs ...uuid.UUID) (*ChatEventsSubscription, error) {
	subscription := c.pubSub.Subscription(subscriberName)

	return &ChatEventsSubscription{subscription}, nil
}

type ChatEventsSubscription struct {
	subscription *Subscription
}

func (s *ChatEventsSubscription) Read(ctx context.Context) (event.Envelope[domain.Event], error) {
	var envelope event.Envelope[domain.Event]

	m, err := s.subscription.Read(ctx)
	if err != nil {
		return envelope, err
	}

	// todo unmarshal RawEnvelope
	err = json.Unmarshal(json.RawMessage(m.Payload), &envelope)
	if err != nil {
		return envelope, err
	}

	return envelope, nil
}

func (s *ChatEventsSubscription) Cancel(ctx context.Context) {
	s.Cancel(ctx)
}

func (s *ChatEventsSubscription) AddChats(ctx context.Context, chatIDs ...uuid.UUID) {
	var channels []string
	for _, cid := range chatIDs {
		channels = append(channels, ChannelNameChatEvents(cid))
	}
	s.subscription.AddChannels(ctx, channels...)
}

func ChannelNameChatEvents(chatID uuid.UUID) string {
	return fmt.Sprintf("chats:%s:events", chatID.String())
}
