package user

import "uuid"

const UserRegisteredEventName = "UserRegistered"

type UserRegistered struct {
	ID uuid.UUID
}

func (m UserRegistered) IsEvent() {}

func (m UserRegistered) Name() string {
	return UserRegisteredEventName
}
