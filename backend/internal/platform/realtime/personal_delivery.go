package realtime

import (
	"context"
	"encoding/json/v2"
	directoryQry "messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	messagingEvent "messenger/messenger/internal/messaging/app/event"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/shared/app/event"
	appQuery "messenger/messenger/internal/shared/app/query"
)

var _ rabbitmq.Handler = (*PersonalDeliveryHandler)(nil)

type PersonalDeliveryHandler struct {
	personalEventsPubSub *pubsub.PersonalEventsPubSub
	getEntry             appQuery.Handler[directoryQry.GetEntry, directoryQry.EntryView]
}

func NewPersonalDeliveryHandler(
	personalEventsPubSub *pubsub.PersonalEventsPubSub,
	getEntry appQuery.Handler[directoryQry.GetEntry, directoryQry.EntryView],
) *PersonalDeliveryHandler {
	return &PersonalDeliveryHandler{
		personalEventsPubSub: personalEventsPubSub,
		getEntry:             getEntry,
	}
}

func (h *PersonalDeliveryHandler) Handle(ctx context.Context, msg rabbitmq.Envelope) error {
	switch msg.Type {
	case messagingEvent.DirectConversationStartedEventName:
		var envlp event.Envelope[messagingEvent.DirectConversationStarted]
		err := json.Unmarshal(msg.Payload, &envlp)
		if err != nil {
			return err
		}

		evnt := envlp.Event

		entr, err := h.getEntry.Handle(ctx, directoryQry.GetEntry{
			ID:   evnt.SenderID,
			Type: entry.TargetTypeUser.String(),
		})
		if err != nil {
			return err
		}

		payload, err := json.Marshal(pubsub.DirectConversationStarted{
			ConversationID:  evnt.ConversationID,
			SenderID:        evnt.SenderID,
			SenderAvatarUrl: entr.AvatarURL,
			SenderName:      entr.DisplayName,
			MessageBody:     evnt.MessageBody,
			MessageSentAt:   evnt.MessageSentAt,
		})
		if err != nil {
			return err
		}

		pubSubMsg := pubsub.Message{
			Type:    pubsub.DirectConversationStartedName,
			Payload: payload,
		}

		err = h.personalEventsPubSub.Publish(ctx, envlp.Event.RecipientID, pubSubMsg)
		if err != nil {
			return err
		}
	}
	return nil
}
