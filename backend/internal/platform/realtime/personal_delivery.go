package realtime

import (
	"context"
	"encoding/json/v2"
	directoryQry "messenger/messenger/internal/directory/app/query"
	directoryView "messenger/messenger/internal/directory/app/view"
	mediaQry "messenger/messenger/internal/media/app/query"
	messagingEvent "messenger/messenger/internal/messaging/app/event"
	"messenger/messenger/internal/platform/pubsub"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/shared/app/event"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

var _ rabbitmq.Handler = (*PersonalDeliveryHandler)(nil)

type PersonalDeliveryHandler struct {
	personalEventsPubSub *pubsub.PersonalEventsPubSub
	getUserProfiles             appQuery.Handler[directoryQry.GetUserProfiles, []directoryView.UserProfile]
	getPublicURLs        appQuery.Handler[mediaQry.GetPublicURLs, map[uuid.UUID]string]
}

func NewPersonalDeliveryHandler(
	personalEventsPubSub *pubsub.PersonalEventsPubSub,
	getUserProfiles appQuery.Handler[directoryQry.GetUserProfiles, []directoryView.UserProfile],
	getPublicURLs appQuery.Handler[mediaQry.GetPublicURLs, map[uuid.UUID]string],
) *PersonalDeliveryHandler {
	return &PersonalDeliveryHandler{
		personalEventsPubSub: personalEventsPubSub,
		getUserProfiles:             getUserProfiles,
		getPublicURLs:        getPublicURLs,
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

		ups, err := h.getUserProfiles.Handle(ctx, directoryQry.GetUserProfiles{UserIDs: []uuid.UUID{evnt.SenderID}})
		if err != nil {
			return err
		}
		userProfile := ups[0]

		var avatarURL *string
		if userProfile.AvatarID != nil {
			urls, err := h.getPublicURLs.Handle(ctx, mediaQry.GetPublicURLs{MediaIDs: []uuid.UUID{*userProfile.AvatarID}})
			if err != nil {
				return err
			}
			if url, ok := urls[*userProfile.AvatarID]; ok {
				avatarURL = &url
			}
		}

		payload, err := json.Marshal(pubsub.DirectConversationStarted{
			ConversationID:  evnt.ConversationID,
			SenderID:        evnt.SenderID,
			SenderAvatarURL: avatarURL,
			SenderName:      userProfile.DisplayName,
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
