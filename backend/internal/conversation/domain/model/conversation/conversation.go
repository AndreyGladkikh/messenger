package conversation

import (
	"bytes"
	"errors"
	"messenger/messenger/internal/conversation/domain/model/participant"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type ConversationType string

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
		id: uuid.New(),
		Type: typ,
	}
}

type DirectConversationParticipantsPair struct {
	participantOne participant.Participant
	participantTwo participant.Participant
}

func NewDirectConversationParticipantsPair(
	participantOne participant.Participant,
	participantTwo participant.Participant,
) (DirectConversationParticipantsPair, error) {
	c := bytes.Compare(participantOne.ID[:], participantTwo.ID[:])
	if c == 0 {
		return DirectConversationParticipantsPair{}, errors.New("direct conversation participants must be different")
	}
	if c == 1 {
		participantOne, participantTwo = participantTwo, participantOne
	}

	return DirectConversationParticipantsPair{
		participantOne: participantOne,
		participantTwo: participantTwo,
	}, nil
}

func (p DirectConversationParticipantsPair) List() [2]participant.Participant {
	return [2]participant.Participant{p.participantOne, p.participantTwo}
}
