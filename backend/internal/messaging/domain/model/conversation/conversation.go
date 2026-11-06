package conversation

import (
	"bytes"
	"errors"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type Type string

func (c Type) String() string {
	return string(c)
}

const (
	TypeDirect Type = "direct"
	TypeGroup  Type = "group"
)

type Conversation struct {
	domain.AggregateRoot

	id        uuid.UUID
	Type      Type
	Title     *string
	AvatarUrl *string
}

func New(
	typ Type,
) *Conversation {
	return &Conversation{
		id:   uuid.New(),
		Type: typ,
	}
}

func Rehydrate(
	id uuid.UUID,
	typ Type,
) *Conversation {
	return &Conversation{
		id:   id,
		Type: typ,
	}
}

func (s *Conversation) ID() uuid.UUID {
	return s.id
}

type DirectConversationParticipantsPair struct {
	participantLow  uuid.UUID
	participantHigh uuid.UUID
}

func NewDirectConversationParticipantsPair(
	participantOne uuid.UUID,
	participantTwo uuid.UUID,
) (DirectConversationParticipantsPair, error) {
	c := bytes.Compare(participantOne[:], participantTwo[:])
	if c == 0 {
		return DirectConversationParticipantsPair{}, errors.New("direct conversation participants must be different")
	}
	if c == 1 {
		participantOne, participantTwo = participantTwo, participantOne
	}

	return DirectConversationParticipantsPair{
		participantLow:  participantOne,
		participantHigh: participantTwo,
	}, nil
}

func (p DirectConversationParticipantsPair) List() [2]uuid.UUID {
	return [2]uuid.UUID{p.participantLow, p.participantHigh}
}

func (p DirectConversationParticipantsPair) ParticipantLow() uuid.UUID {
	return p.participantLow
}

func (p DirectConversationParticipantsPair) ParticipantHigh() uuid.UUID {
	return p.participantHigh
}
