package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"messenger/messenger/internal/platform/event"

	"github.com/google/uuid"
)

type ChatEventsPubSub struct {
	*PubSub
}

func NewChatEventsPubSub(
	pubsub *PubSub,
) *ChatEventsPubSub {
	return &ChatEventsPubSub{
		PubSub: pubsub,
	}
}

func (c *ChatEventsPubSub) Publish(ctx context.Context, chatID uuid.UUID, e event.Envelope) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to publish message to chat events pubsub: %w", err)
	}

	c.PubSub.Publish(ctx, ChannelNameChatEvents(chatID), payload)

	return nil
}

func (c *ChatEventsPubSub) Subscribe(ctx context.Context, subscriberName string, chatIDs ...uuid.UUID) {

}

func ChannelNameChatEvents(chatID uuid.UUID) string {
	return fmt.Sprintf("chats:%s:events", chatID.String())
}