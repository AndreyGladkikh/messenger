package message_sent_notifier

import (
	"context"
	"messenger/messenger/internal/messaging/domain/message"
	"messenger/messenger/internal/platform/redis"
)

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
	_, err := n.pubSubHub.Publish(ctx, redis.ChannelNameChatEvents(msg.ChatID), msg)
	if err != nil {
		return err
	}
	return nil
}
