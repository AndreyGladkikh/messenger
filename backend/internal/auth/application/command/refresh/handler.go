package refresh

import (
	"context"
	"errors"
	"messenger/messenger/internal/auth/application/token"
	"messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/auth/domain/session"
	sharedDomain "messenger/messenger/internal/shared/domain"
)

type Handler struct {
	sessionRepo  session.Repository
	tokenService token.Service
}

func NewHandler(
	sessionRepo session.Repository,
	tokenService token.Service,
) *Handler {
	return &Handler{
		sessionRepo:  sessionRepo,
		tokenService: tokenService,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) (any, error) {
	ses, err := h.sessionRepo.GetByRefreshTokenHash(ctx, h.tokenService.HashRefreshToken(cmd.RefreshToken))
	if err != nil && !errors.Is(err, sharedDomain.ErrNotFound) {
		return nil, err
	}
	if errors.Is(err, sharedDomain.ErrNotFound) {
		return nil, domain.ErrUnauthorized
	}

	refreshToken := h.tokenService.GenerateRefreshToken()
	refreshTokenHash := h.tokenService.HashRefreshToken(refreshToken)

	ses.Refresh(refreshTokenHash)

	if err = h.sessionRepo.Save(ctx, ses); err != nil {
		return nil, err
	}

	accessToken, err := h.tokenService.GenerateAccessToken(ses)
	if err != nil {
		return nil, err
	}

	return &Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
