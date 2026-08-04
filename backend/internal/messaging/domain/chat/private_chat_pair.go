package chat

import (
	"bytes"
	"errors"

	"github.com/google/uuid"
)

type PrivateChatPair struct {
	firstParticipantID  uuid.UUID
	secondParticipantID uuid.UUID
}

func NewPrivateChatPair(firstParticipantID, secondParticipantID uuid.UUID) (*PrivateChatPair, error) {
	c := bytes.Compare(firstParticipantID[:], secondParticipantID[:])
	if c == 0 {
		return nil, errors.New("participants must be different")
	}
	if c == 1 {
		firstParticipantID, secondParticipantID = secondParticipantID, firstParticipantID
	}

	return &PrivateChatPair{
		firstParticipantID:  firstParticipantID,
		secondParticipantID: secondParticipantID,
	}, nil
}

func (p *PrivateChatPair) FirstParticipantID() uuid.UUID {
	return p.firstParticipantID
}

func (p *PrivateChatPair) SecondParticipantID() uuid.UUID {
	return p.secondParticipantID
}
