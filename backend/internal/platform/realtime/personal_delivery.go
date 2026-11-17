package realtime

import (
	"context"
	"encoding/json/v2"
	directoryQry "messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/media/app/query"
	messagingEvent "messenger/messenger/internal/messaging/app/event"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/query_handler"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/shared/app/event"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

var _ rabbitmq.Handler = (*PersonalDeliveryHandler)(nil)

type PersonalDeliveryHandler struct {
	personalEventsPubSub *pubsub.PersonalEventsPubSub
	getEntry             appQuery.Handler[directoryQry.GetEntry, directoryQry.EntryView]
	queries *query_handler.Registry
}

func NewPersonalDeliveryHandler(
	personalEventsPubSub *pubsub.PersonalEventsPubSub,
	getEntry appQuery.Handler[directoryQry.GetEntry, directoryQry.EntryView],
	queries *query_handler.Registry,
) *PersonalDeliveryHandler {
	return &PersonalDeliveryHandler{
		personalEventsPubSub: personalEventsPubSub,
		getEntry:             getEntry,
		queries: queries,
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

		var avatarURL *string
		if entr.AvatarID != nil {
			urls, err := h.queries.GetMediaPublicURLsHandler.Handle(ctx, query.GetPublicURLs{MediaIDs: []uuid.UUID{*entr.AvatarID}})
			if err != nil {
				return err
			}
			if url, ok := urls[*entr.AvatarID]; ok {
				avatarURL = &url
			}
		}

		payload, err := json.Marshal(pubsub.DirectConversationStarted{
			ConversationID:  evnt.ConversationID,
			SenderID:        evnt.SenderID,
			SenderAvatarURL: avatarURL,
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
