package user

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
	Name         string
}

func Register(
	login string,
	passwordHash string,
) *User {
	return &User{
		ID:           uuid.New(),
		Login:        login,
		PasswordHash: passwordHash,
	}
}

func Rehydrate(
	id uuid.UUID,
	login string,
	passwordHash string,
	name string,
) *User {
	return &User{
		ID:           id,
		Login:        login,
		PasswordHash: passwordHash,
		Name:         name,
	}
}
