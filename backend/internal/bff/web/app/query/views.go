package query

import (
	"time"
	"uuid"
)

type PeerItem struct {
	ID          uuid.UUID
	Type        string
	AvatarURL   *string
	Title       string
	Preview     string
	PreviewTime time.Time
	UnreadCount int
}
