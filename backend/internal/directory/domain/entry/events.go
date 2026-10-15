package entry

import "uuid"

const EntryRegisteredEventName = "EntryRegistered"

type EntryRegistered struct {
	ID          uuid.UUID
	TargetID    uuid.UUID
	TargetType  TargetType
	DisplayName string
	Handle      string
}

func (e EntryRegistered) IsEvent() {}

func (e EntryRegistered) Name() string {
	return EntryRegisteredEventName
}
