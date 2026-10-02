package chat_participant

import "uuid"

type ParticipantAddedToChat struct {
	ChatID        uuid.UUID
	ParticipantID uuid.UUID
}

func (c ParticipantAddedToChat) IsEvent() {}

func (p ParticipantAddedToChat) Name() string {
	return "ParticipantAddedToChat"
}
