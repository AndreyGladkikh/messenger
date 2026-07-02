package chat

import "messenger/messenger/internal/domain/chat_participant"

type Chat struct {
	id string
	typ ChatType
	name string
	participants []chat_participant.ChatParticipant
}
