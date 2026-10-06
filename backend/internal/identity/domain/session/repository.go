package session

import (
	"context"
)

type Repository interface {
	GetByRefreshTokenHash(context.Context, string) (*Session, error)
	Save(context.Context, *Session) error
}
