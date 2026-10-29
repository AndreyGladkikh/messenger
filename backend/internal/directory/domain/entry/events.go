package entry

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

const EntryRegisteredEventName = "EntryRegistered"

var _ domain.Event = EntryRegistered{}

type EntryRegistered struct {
	ID          uuid.UUID
	TargetID    uuid.UUID
	TargetType  TargetType
	DisplayName string
	Handle      string
}

func (e EntryRegistered) Name() string {
	return EntryRegisteredEventName
}
