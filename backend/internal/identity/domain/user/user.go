package user

import "uuid"

type User struct {
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

func (u *User) ID() uuid.UUID {
	return u.id
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
