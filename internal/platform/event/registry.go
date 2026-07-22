package event

import "messenger/messenger/internal/application/event/message_sent"

type HandlerRegistry struct {
	NotifyChatParticipantsHandler *message_sent.NotifyChatParticipantsHandler
}