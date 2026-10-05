package account

import (
	"context"
	"messenger/messenger/internal/identity/domain/user"
)

type Repository interface {
	SaveUser(context.Context, *user.User) error
}
