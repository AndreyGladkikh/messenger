package port

import (
	"context"
	"uuid"
)

const (
	DirectoryEntryTypeUser    = "user"
	DirectoryEntryTypeBot     = "bot"
	DirectoryEntryTypeGroup   = "group"
	DirectoryEntryTypeChannel = "channel"
)

type EntryRef struct {
	ID   uuid.UUID
	Type string
}

type EntryView struct {
	ID          uuid.UUID
	Type        string
	DisplayName string
	AvatarURL   *string
}

type DirectoryReader interface {
	GetEntries(ctx context.Context, refs []EntryRef) ([]EntryView, error)
}
