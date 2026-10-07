package user

import (
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type User struct {
	domain.AggregateRoot
	
	id           uuid.UUID
	Login        string
	PasswordHash string
}

func New(
	login string,
	passwordHash string,
) *User {
	return &User{
		id:           uuid.New(),
		Login:        login,
		PasswordHash: passwordHash,
	}
}

func Rehydrate(
	id uuid.UUID,
	login string,
	passwordHash string,
) *User {
	return &User{
		id:           id,
		Login:        login,
		PasswordHash: passwordHash,
	}
}

func (u *User) ID() uuid.UUID {
	return u.id
}
