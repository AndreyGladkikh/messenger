package login

import (
	"context"
	"errors"
	"messenger/messenger/internal/auth/application/password"
	"messenger/messenger/internal/auth/application/token"
	"messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/auth/domain/session"
	"messenger/messenger/internal/auth/domain/user"
)

type Handler struct {
	userRepo     user.Repository
	sessionRepo  session.Repository
	passHasher   password.Hasher
	tokenService token.Service
}

func NewHandler(
	userRepo user.Repository,
	sessionRepo session.Repository,
	passHasher password.Hasher,
	tokenService token.Service,
) *Handler {
	return &Handler{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		passHasher:   passHasher,
		tokenService: tokenService,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) (any, error) {
	u, err := h.userRepo.GetByLogin(ctx, cmd.Login)
	if err != nil && !errors.Is(err, user.ErrNotFound) {
		return nil, err
	}
	if errors.Is(err, user.ErrNotFound) {
		return nil, domain.ErrInvalidCredantials
	}

	eq, err := h.passHasher.Compare(cmd.Password, u.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !eq {
		return nil, domain.ErrInvalidCredantials
	}

	refreshToken := h.tokenService.GenerateRefreshToken()
	refreshTokenHash := h.tokenService.HashRefreshToken(refreshToken)

	session := session.Create(
		u.ID,
		refreshTokenHash,
		cmd.UserAgent,
		cmd.IP,
	)
	if err = h.sessionRepo.Add(ctx, session); err != nil {
		return nil, err
	}

	accessToken, err := h.tokenService.GenerateAccessToken(session)
	if err != nil {
		return nil, err
	}

	return &Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
