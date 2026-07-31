package session

import (
	"net/netip"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	CreatedAt        time.Time
	ExpiresAt        time.Time
	RevokedAt        time.Time
	LastUsedAt       time.Time
	UserAgent        string
	IP               netip.Addr
}

func Create(
	userID uuid.UUID,
	refreshTokenHash string,
	userAgent string,
	ip netip.Addr,
) *Session {
	createdAt := time.Now()
	expiresAt := createdAt.AddDate(0, 1, 0)

	return &Session{
		ID:               uuid.New(),
		UserID:           userID,
		CreatedAt:        createdAt,
		ExpiresAt:        expiresAt,
		RefreshTokenHash: refreshTokenHash,
		UserAgent:        userAgent,
		IP:               ip,
	}
}

func Rehydrate(
	id uuid.UUID,
	userID uuid.UUID,
	refreshTokenHash string,
	createdAt time.Time,
	expiresAt time.Time,
	revokedAt time.Time,
	lastUsedAt time.Time,
	userAgent string,
	ip netip.Addr,
) *Session {
	return &Session{
		ID:               id,
		UserID:           userID,
		RefreshTokenHash: refreshTokenHash,
		CreatedAt:        createdAt,
		ExpiresAt:        expiresAt,
		RevokedAt:        revokedAt,
		LastUsedAt:       lastUsedAt,
		UserAgent:        userAgent,
		IP:               ip,
	}
}

func (s *Session) IsRevoked() bool {
	return !s.RevokedAt.IsZero()
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) Revoke() {
	if s.IsRevoked() {
		return
	}
	s.RevokedAt = time.Now()
}
