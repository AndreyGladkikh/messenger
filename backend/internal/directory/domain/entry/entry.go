package entry

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type Entry struct {
	domain.AggregateRoot

	id          uuid.UUID
	TargetID    uuid.UUID
	TargetType  TargetType
	DisplayName string
	Handle      string
}

func New(
	targetID uuid.UUID,
	targetType TargetType,
	displayName string,
	handle string,
) *Entry {
	e := &Entry{
		id:          uuid.New(),
		TargetID:    targetID,
		TargetType:  targetType,
		DisplayName: displayName,
		Handle:      handle,
	}
	e.AddEvent(EntryRegistered{
		ID:          e.ID(),
		TargetID:    targetID,
		TargetType:  targetType,
		DisplayName: displayName,
		Handle:      handle,
	})
	return e
}

func (e *Entry) ID() uuid.UUID {
	return e.id
}
