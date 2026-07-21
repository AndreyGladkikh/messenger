package chat_participant

type ParticipantAddedToChat struct {
	ChatID string
	ParticipantID string
}

func (p ParticipantAddedToChat) Name() string {
	return "ParticipantAddedToChat"
}