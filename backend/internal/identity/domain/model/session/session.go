package session

import (
	"messenger/messenger/internal/shared/domain"
	"net/netip"
	"time"

	"uuid"
)

const TTL = 30 * 24 * time.Hour

type Session struct {
	domain.AggregateRoot

	id               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	CreatedAt        time.Time
	ExpiresAt        time.Time
	RevokedAt        *time.Time
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
	s := &Session{
		id:         uuid.New(),
		UserID:     userID,
		CreatedAt:  time.Now(),
		UserAgent:  userAgent,
		IP:         ip,
		LastUsedAt: time.Now(),
	}
	s.Refresh(refreshTokenHash)
	return s
}

func Rehydrate(
	id uuid.UUID,
	userID uuid.UUID,
	refreshTokenHash string,
	createdAt time.Time,
	expiresAt time.Time,
	revokedAt *time.Time,
	lastUsedAt time.Time,
	userAgent string,
	ip netip.Addr,
) *Session {
	return &Session{
		id:               id,
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

func (s *Session) ID() uuid.UUID {
	return s.id
}

func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) Revoke() {
	if s.IsRevoked() {
		return
	}
	now := time.Now()
	s.RevokedAt = &now
}

func (s *Session) Refresh(tokenHash string) {
	s.RefreshTokenHash = tokenHash
	s.ExpiresAt = time.Now().Add(TTL)
}
