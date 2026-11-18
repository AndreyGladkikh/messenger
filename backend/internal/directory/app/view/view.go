package view

import (
	"time"
	"uuid"
)

type Entry struct {
	ID          uuid.UUID
	Type        string
	Handle      string
	DisplayName string
	AvatarID    *uuid.UUID
}

type UserProfile struct {
	UserID uuid.UUID
	Handle      *string
	DisplayName string
	AvatarID    *uuid.UUID
	Bio         *string
	Birthday    *time.Time
}
