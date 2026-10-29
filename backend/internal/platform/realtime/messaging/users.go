package messaging

import (
	"context"
	"encoding/json/v2"
	cnvEvent "messenger/messenger/internal/conversations/app/event"
	drctrQry "messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/shared/app/event"
)

var _ rabbitmq.Consumer = (*UserDeliveryConsumer)(nil)

type UserDeliveryConsumer struct {
	UserPubSub *pubsub.UserEventsPubSub
	QueryBus   *querybus.Bus
}

func (h *UserDeliveryConsumer) HandleMessage(ctx context.Context, msg rabbitmq.Message) error {
	switch msg.Type {
	case cnvEvent.DirectConversationWithFirstMessageCreatedEventName:
		var envlp event.Envelope[cnvEvent.DirectConversationWithFirstMessageCreated]
		err := json.Unmarshal(msg.Payload, envlp)
		if err != nil {
			return err
		}

		evnt := envlp.Event

		entr, err := h.QueryBus.Dispatch[drctrQry.EntryView](ctx, drctrQry.GetEntry{
			ID:   evnt.SenderID,
			Type: entry.TargetTypeUser.String(),
		})
		if err != nil {
			return err
		}

		payload, err := json.Marshal(pubsub.DirectConversationCreated{
			ConversationID:  evnt.ConversationID,
			SenderAvatarUrl: *entr.AvatarURL,
			SenderName:      entr.DisplayName,
			MessageBody:     evnt.MessageBody,
			MessageSentAt:   envlp.OccurredAt,
		})
		if err != nil {
			return err
		}

		pubSubMsg := pubsub.Message{
			Type:    pubsub.DirectConversationCreatedName,
			Payload: payload,
		}

		h.UserPubSub.Publish(ctx, envlp.Event.RecipientID, pubSubMsg)
	}
	return nil
}
