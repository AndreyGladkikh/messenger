package user

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

const UserRegisteredEventName = "identity.user.Registered"

var _ domain.Event = UserRegistered{}

type UserRegistered struct {
	ID uuid.UUID
}

func (m UserRegistered) Name() string {
	return UserRegisteredEventName
}
