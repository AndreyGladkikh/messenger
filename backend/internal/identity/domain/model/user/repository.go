package user

import (
	"context"
)

type Repository interface {
	ExistsByLogin(context.Context, string) (bool, error)
	Save(context.Context, *User) error
	GetByLogin(context.Context, string) (*User, error)
}
