package view

import (
	"time"
	"uuid"
)

type PeerItem struct {
	ID          uuid.UUID
	Type        string
	AvatarURL   *string
	Title       string
	Handle      *string
	Preview     *string
	PreviewTime *time.Time
	UnreadCount int
}

type UserProfile struct {
	UserID      uuid.UUID
	Handle      *string
	DisplayName string
	Bio         *string
	Birthday    *time.Time
	AvatarURL   *string
}
