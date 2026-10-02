package participant

import "uuid"

type ParticipantType string

const (
	ParticipantTypeUser ParticipantType = "user"
	ParticipantTypeBot ParticipantType = "bot"
)

type Participant struct {
	ID   uuid.UUID
	Type ParticipantType
}
