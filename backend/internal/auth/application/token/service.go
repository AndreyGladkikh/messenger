package token

import (
	"messenger/messenger/internal/auth/domain/session"
	"time"
)

const AccessTokenTTL = 15 * time.Minute

type Service interface {
	GenerateAccessToken(*session.Session) (string, error)
	GenerateRefreshToken() string
	HashRefreshToken(token string) string
}
