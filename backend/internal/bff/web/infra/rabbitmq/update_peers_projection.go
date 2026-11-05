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
		// var envlp appEvent.Envelope[cnvEvent.DirectConversationStarted]
		// err := json.Unmarshal(msg.Payload, &envlp)
		// if err != nil {
		// 	return err
		// }

		// evnt := envlp.Event

		// entr, err := h.queryBus.Dispatch[drctrQry.EntryView](ctx, drctrQry.GetEntry{
		// 	ID:   evnt.SenderID,
		// 	Type: entry.TargetTypeUser.String(),
		// })
		// if err != nil {
		// 	return err
		// }

		// payload, err := json.Marshal(pubsub.DirectConversationStarted{
		// 	ConversationID:  evnt.ConversationID,
		// 	SenderAvatarUrl: entr.AvatarURL,
		// 	SenderName:      entr.DisplayName,
		// 	MessageBody:     evnt.MessageBody,
		// 	MessageSentAt:   envlp.OccurredAt,
		// })
		// if err != nil {
		// 	return err
		// }

		// pubSubMsg := pubsub.Message{
		// 	Type:    pubsub.DirectConversationStartedName,
		// 	Payload: payload,
		// }

		// err = h.userPubSub.Publish(ctx, envlp.Event.RecipientID, pubSubMsg)
		// if err != nil {
		// 	return err
		// }
	case message.MessageSentEventName:

	}

	return nil
}
