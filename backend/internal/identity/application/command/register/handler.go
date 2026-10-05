package register

import (
	"context"
	"fmt"
	"messenger/messenger/internal/identity/application/password"
	"messenger/messenger/internal/identity/application/token"
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

	exists, err := h.userRepo.ExistsByLogin(ctx, cmd.Login)
	if err != nil {
		return zero, err
	}
	if exists {
		return zero, user.ErrLoginAlreadyExists
	}

	passHash, err := h.passHasher.Hash(cmd.Password)
	if err != nil {
		return zero, fmt.Errorf("failed to hash password: %w", err)
	}

	user := user.Register(
		cmd.Login,
		passHash,
	)

	err = h.userRepo.Add(ctx, user)
	if err != nil {
		return zero, err
	}

	refreshToken := h.tokenService.GenerateRefreshToken()
	refreshTokenHash := h.tokenService.HashRefreshToken(refreshToken)

	session := session.Create(
		user.ID,
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
