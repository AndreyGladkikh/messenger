package chat_participant

import "messenger/messenger/internal/domain"

type ChatParticipant struct {
	domain.BaseAggregate

	id string
	chatID    string
	participantID string
	role Role
}

func (p *ChatParticipant) ID() string {
	return p.id
}

func (p *ChatParticipant) ChatID() string {
	return p.chatID
}

func (p *ChatParticipant) ParticipantID() string {
	return p.participantID
}

func (p *ChatParticipant) Role() Role {
	return p.role
}