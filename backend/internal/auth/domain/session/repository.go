package session

import (
	"context"
)

type Repository interface {
	Add(context.Context, *Session) error
	GetByRefreshTokenHash(context.Context, string) (*Session, error)
	Save(context.Context, *Session) error
}
