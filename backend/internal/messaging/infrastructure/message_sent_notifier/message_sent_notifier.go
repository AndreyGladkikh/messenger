package message_sent_notifier

import (
	"context"
	"fmt"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/platform/redis"

	"github.com/google/uuid"
)

func chName(chatID uuid.UUID) string {
	return fmt.Sprintf("chats:%s:events", chatID.String())
}

type MessageSentNotifier struct {
	pubSubHub *redis.RedisPubSubHub
}

func NewMessageSentNotifier(
	pubSubHub *redis.RedisPubSubHub,
) *MessageSentNotifier {
	return &MessageSentNotifier{
		pubSubHub: pubSubHub,
	}
}

func (n *MessageSentNotifier) Notify(ctx context.Context, msg message.MessageSent) error {
	_, err := n.pubSubHub.Publish(ctx, chName(msg.ChatID), msg)
	if err != nil {
		return err
	}
	return nil
}