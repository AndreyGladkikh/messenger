package conversation

import (
	"bytes"
	"errors"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type ConversationType string

func (c ConversationType) String() string {
	return string(c)
}

const (
	ConversationTypeDirect ConversationType = "direct"
	ConversationTypeGroup  ConversationType = "group"
)

var _ domain.Aggregate = (*Conversation)(nil)

type Conversation struct {
	domain.AggregateRoot

	id   uuid.UUID
	Type ConversationType
}

func (s *Conversation) ID() uuid.UUID {
	return s.id
}

func New(
	typ ConversationType,
) *Conversation {
	return &Conversation{
		id:   uuid.New(),
		Type: typ,
	}
}

func Rehydrate(
	id uuid.UUID,
	typ ConversationType,
) *Conversation {
	return &Conversation{
		id:   id,
		Type: typ,
	}
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
