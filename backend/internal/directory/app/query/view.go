package query

import "uuid"

type EntryView struct {
	ID          uuid.UUID
	Type        string
	DisplayName string
	AvatarURL   *string
}
