package notifier

import (
	"context"
	"messenger/messenger/internal/messaging/domain/message"
)

type MessageSentNotifier interface {
	Notify(context.Context, message.MessageSent) error
}