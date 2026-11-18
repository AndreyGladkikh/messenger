package port

import (
	"context"
	"time"
	"uuid"
)

const (
	DirectoryEntryTypeUser    = "user"
	DirectoryEntryTypeBot     = "bot"
	DirectoryEntryTypeGroup   = "group"
	DirectoryEntryTypeChannel = "channel"
)

type DirectoryReader interface {
	GetUserProfiles(ctx context.Context, ids []uuid.UUID) ([]UserProfileView, error)
	Search(ctx context.Context, str string) ([]SearchResultItem, error)
}

type UserProfileView struct {
	UserID      uuid.UUID
	Handle      *string
	DisplayName string
	AvatarID    *uuid.UUID
	Bio         *string
	Birthday    *time.Time
}

type SearchResultItem struct {
	ID          uuid.UUID
	Type        string
	Handle      *string
	DisplayName string
	AvatarID    *uuid.UUID
}
