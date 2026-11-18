package user_profile

import (
	"messenger/messenger/internal/shared/domain"
	"time"

	"uuid"
)

const CreatedEventName = "directory.user_profile.Created.v1"

var _ domain.Event = Created{}

type Created struct {
	UserID      uuid.UUID
	DisplayName string
	AvatarID    *uuid.UUID
	Bio         *string
	Birthday    *time.Time
}

func (c Created) Name() string {
	return CreatedEventName
}

const UpdatedEventName = "directory.user_profile.Updated.v1"

var _ domain.Event = Updated{}

type Updated struct {
	UserID      uuid.UUID
	DisplayName string
	AvatarID    *uuid.UUID
	Bio         *string
	Birthday    *time.Time
}

func (c Updated) Name() string {
	return CreatedEventName
}
