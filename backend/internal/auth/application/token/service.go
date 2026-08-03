package token

import (
	"messenger/messenger/internal/auth/domain/session"
)

type Service interface {
	GenerateAccessToken(*session.Session) (string, error)
	GenerateRefreshToken() string
	HashRefreshToken(token string) string
}
