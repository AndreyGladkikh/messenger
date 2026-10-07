package repository

import (
	"context"
	"database/sql"
	"errors"
	"messenger/messenger/internal/identity/domain/session"
	"messenger/messenger/internal/platform/repository"
	"messenger/messenger/internal/platform/sqlc"
	sharedDomain "messenger/messenger/internal/shared/domain"
)

type SessionRepository struct {
	*repository.Repository
}

func NewSessionRepository(repo *repository.Repository) *SessionRepository {
	return &SessionRepository{
		Repository: repo,
	}
}

func (r *SessionRepository) GetByRefreshTokenHash(ctx context.Context, token string) (*session.Session, error) {
	s, err := r.Queries(ctx).GetSessionByRefreshTokenHash(ctx, token)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sharedDomain.ErrNotFound
	}

	return session.Rehydrate(
		s.ID,
		s.UserID,
		s.RefreshTokenHash,
		s.CreatedAt,
		s.ExpiresAt,
		s.RevokedAt,
		s.LastUsedAt,
		s.UserAgent,
		s.Ip,
	), nil
}

func (r *SessionRepository) Save(ctx context.Context, s *session.Session) error {
	return r.Queries(ctx).UpsertSession(ctx, sqlc.UpsertSessionParams{
		ID:               s.ID,
		UserID:           s.UserID,
		UserAgent:        s.UserAgent,
		Ip:               s.IP,
		RefreshTokenHash: s.RefreshTokenHash,
		ExpiresAt:        s.ExpiresAt,
		RevokedAt:        s.RevokedAt,
		LastUsedAt:       s.LastUsedAt,
	})
}
