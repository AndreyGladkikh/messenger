package chat_participant

type ParticipantAddedToChat struct {
	ChatID string
	ParticipantID string
}

func (c ParticipantAddedToChat) IsEvent() {}

func (p ParticipantAddedToChat) Name() string {
	return "ParticipantAddedToChat"
}