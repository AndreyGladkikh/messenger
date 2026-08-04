package user

import (
	"context"
)

type Repository interface {
	ExistsByLogin(context.Context, string) (bool, error)
	Add(context.Context, *User) error
	GetByLogin(context.Context, string) (*User, error)
}
