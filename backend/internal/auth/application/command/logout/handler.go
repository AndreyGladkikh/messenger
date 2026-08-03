package logout

import (
	"context"
	"errors"
	"messenger/messenger/internal/auth/application/password"
	"messenger/messenger/internal/auth/application/token"
	"messenger/messenger/internal/auth/domain/session"
	"messenger/messenger/internal/auth/domain/user"
	sharedDomain "messenger/messenger/internal/shared/domain"
)

type Handler struct {
	userRepo     user.Repository
	sessionRepo  session.Repository
	passHasher   password.Hasher
	tokenService token.Service
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) (any, error) {
	session, err := h.sessionRepo.GetByRefreshTokenHash(ctx, h.tokenService.HashRefreshToken(cmd.RefreshToken))
	if err != nil && !errors.Is(err, sharedDomain.ErrNotFound) {
		return nil, err
	}
	if errors.Is(err, sharedDomain.ErrNotFound) || session.IsExpired() || session.IsRevoked() {
		return nil, nil
	}

	session.Revoke()

	err = h.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
