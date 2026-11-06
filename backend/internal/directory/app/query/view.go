package query

import "uuid"

type EntryView struct {
	ID          uuid.UUID
	Type        string
	Handle string
	DisplayName string
	AvatarURL   *string
}
