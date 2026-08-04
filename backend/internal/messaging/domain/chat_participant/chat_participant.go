package chat_participant

import (
	sharedDomain "messenger/messenger/internal/shared/domain"

	"github.com/google/uuid"
)

type ChatParticipant struct {
	sharedDomain.BaseAggregate

	id            uuid.UUID
	chatID        uuid.UUID
	participantID uuid.UUID
	role          Role
}

func AddToChat(
	chatID uuid.UUID,
	participantID uuid.UUID,
	role Role,
) *ChatParticipant {
	p := &ChatParticipant{
		id:            uuid.New(),
		chatID:        chatID,
		participantID: participantID,
		role:          role,
	}

	p.AddEvent(ParticipantAddedToChat{
		ChatID:        chatID,
		ParticipantID: participantID,
	})

	return p
}

func (p *ChatParticipant) ID() uuid.UUID {
	return p.id
}

func (p *ChatParticipant) ChatID() uuid.UUID {
	return p.chatID
}

func (p *ChatParticipant) ParticipantID() uuid.UUID {
	return p.participantID
}

func (p *ChatParticipant) Role() Role {
	return p.role
}
