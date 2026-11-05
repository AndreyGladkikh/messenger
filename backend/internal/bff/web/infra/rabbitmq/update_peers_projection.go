package rabbitmq

import (
	"context"
	// "encoding/json/v2"
	cnvEvent "messenger/messenger/internal/conversations/app/event"
	"messenger/messenger/internal/conversations/domain/model/message"
	"messenger/messenger/internal/platform/rabbitmq"
	// appEvent "messenger/messenger/internal/shared/app/event"
)

type PeersProjectionUpdater struct {
}

func (h *PeersProjectionUpdater) Handle(ctx context.Context, msg rabbitmq.Envelope) error {
	switch msg.Type {
	case cnvEvent.DirectConversationStartedEventName:
		return nil
	case message.MessageSentEventName:

	}

	return nil
}
