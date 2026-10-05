package login

import (
	"context"
	"errors"
	"messenger/messenger/internal/identity/application/password"
	"messenger/messenger/internal/identity/application/token"
	"messenger/messenger/internal/identity/domain"
	"messenger/messenger/internal/identity/domain/session"
	"messenger/messenger/internal/identity/domain/user"
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

func (h *Handler) Handle(ctx context.Context, cmd Command) (Response, error) {
	var zero Response

	u, err := h.userRepo.GetByLogin(ctx, cmd.Login)
	if err != nil && !errors.Is(err, user.ErrNotFound) {
		return zero, err
	}
	if errors.Is(err, user.ErrNotFound) {
		return zero, domain.ErrInvalidCredantials
	}

	eq, err := h.passHasher.Compare(cmd.Password, u.PasswordHash)
	if err != nil {
		return zero, err
	}
	if !eq {
		return zero, domain.ErrInvalidCredantials
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
		return zero, err
	}

	accessToken, err := h.tokenService.GenerateAccessToken(session)
	if err != nil {
		return zero, err
	}

	return Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
