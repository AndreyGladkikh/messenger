package membership

import (
	"messenger/messenger/internal/conversation/domain/model/participant"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

var _ domain.Aggregate = (*Membership)(nil)

type Membership struct {
	domain.AggregateRoot

	id             uuid.UUID
	ConversationID uuid.UUID
	Participant    participant.Participant
	Role           Role
}

func (m *Membership) ID() uuid.UUID {
	return m.id
}

func New(
	conversationID uuid.UUID,
	participant participant.Participant,
	role Role,
) *Membership {
	return &Membership{
		id:             uuid.New(),
		ConversationID: conversationID,
		Participant:    participant,
		Role:           role,
	}
}
