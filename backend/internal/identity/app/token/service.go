package token

import (
	"messenger/messenger/internal/identity/domain/model/session"
	"time"
)

const AccessTokenTTL = 15 * time.Minute

type Service interface {
	GenerateAccessToken(*session.Session) (string, error)
	GenerateRefreshToken() string
	HashRefreshToken(token string) string
}
