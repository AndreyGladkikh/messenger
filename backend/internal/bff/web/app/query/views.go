package query

import (
	"time"
	"uuid"
)

type PeerItem struct {
	ID          uuid.UUID
	Type        string
	AvatarUrl   string
	Title       string
	Preview     string
	PreviewTime time.Time
	UnreadCount int
}
