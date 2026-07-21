package chat_participant

import "messenger/messenger/internal/domain"

type ChatParticipant struct {
	domain.BaseAggregate

	id string
	chatID    string
	participantID string
	role Role
}

func AddToChat(
	id string,
	chatID    string,
	participantID string,
	role Role,
) *ChatParticipant {
	p := &ChatParticipant{
		id: id,
		chatID: chatID,
		participantID: participantID,
		role: role,
	}

	p.AddEvent(ParticipantAddedToChat{
		ChatID: chatID,
		ParticipantID: participantID,
	})

	return p
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