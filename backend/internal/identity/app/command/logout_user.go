package command

import (
	"context"
	"errors"
	"messenger/messenger/internal/identity/app/password"
	"messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/session"
	"messenger/messenger/internal/identity/domain/user"
	"messenger/messenger/internal/shared/domain"
)

const LogoutUserCommandName = "LogoutUser"

type LogoutUser struct {
	RefreshToken string
}

func (c LogoutUser) IsCommand() {}

func (c LogoutUser) Name() string {
	return LogoutUserCommandName
}

type LogoutUserHandler struct {
	userRepo     user.Repository
	sessionRepo  session.Repository
	passHasher   password.Hasher
	tokenService token.Service
}

func (h *LogoutUserHandler) Handle(ctx context.Context, cmd LogoutUser) (any, error) {
	session, err := h.sessionRepo.GetByRefreshTokenHash(ctx, h.tokenService.HashRefreshToken(cmd.RefreshToken))
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	if errors.Is(err, domain.ErrNotFound) || session.IsExpired() || session.IsRevoked() {
		return nil, nil
	}

	session.Revoke()

	err = h.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
