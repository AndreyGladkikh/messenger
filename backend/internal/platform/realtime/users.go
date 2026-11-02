package realtime

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

var _ rabbitmq.Handler = (*PersonalDeliveryHandler)(nil)

type PersonalDeliveryHandler struct {
	userPubSub *pubsub.PersonalEventsPubSub
	queryBus   *querybus.Bus
}

func NewPersonalDeliveryHandler(
	userPubSub *pubsub.PersonalEventsPubSub,
	queryBus *querybus.Bus,
) *PersonalDeliveryHandler {
	return &PersonalDeliveryHandler{
		userPubSub: userPubSub,
		queryBus:   queryBus,
	}
}

func (h *PersonalDeliveryHandler) Handle(ctx context.Context, msg rabbitmq.Envelope) error {
	switch msg.Type {
	case cnvEvent.DirectConversationStartedEventName:
		var envlp event.Envelope[cnvEvent.DirectConversationStarted]
		err := json.Unmarshal(msg.Payload, &envlp)
		if err != nil {
			return err
		}

		evnt := envlp.Event

		entr, err := h.queryBus.Dispatch[drctrQry.EntryView](ctx, drctrQry.GetEntry{
			ID:   evnt.SenderID,
			Type: entry.TargetTypeUser.String(),
		})
		if err != nil {
			return err
		}

		payload, err := json.Marshal(pubsub.DirectConversationStarted{
			ConversationID:  evnt.ConversationID,
			SenderAvatarUrl: entr.AvatarURL,
			SenderName:      entr.DisplayName,
			MessageBody:     evnt.MessageBody,
			MessageSentAt:   envlp.OccurredAt,
		})
		if err != nil {
			return err
		}

		pubSubMsg := pubsub.Message{
			Type:    pubsub.DirectConversationStartedName,
			Payload: payload,
		}

		err = h.userPubSub.Publish(ctx, envlp.Event.RecipientID, pubSubMsg)
		if err != nil {
			return err
		}
	}
	return nil
}
