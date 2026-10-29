package command

import (
	"context"
	"errors"
	"messenger/messenger/internal/identity/app"
	"messenger/messenger/internal/identity/app/password"
	"messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/identity/domain/model/user"
	"messenger/messenger/internal/shared/app/apperr"
	"net/netip"
)

const LoginUserCommandName = "LoginUser"

type LoginUser struct {
	Login     string
	Password  string
	UserAgent string
	IP        netip.Addr
}

func (c LoginUser) IsCommand() {}

func (c LoginUser) Name() string {
	return LoginUserCommandName
}

type LoginUserResponse struct {
	AccessToken  string
	RefreshToken string
}

type LoginUserHandler struct {
	userRepo     user.Repository
	sessionRepo  session.Repository
	passHasher   password.Hasher
	tokenService token.Service
}

func NewLoginUserHandler(
	userRepo user.Repository,
	sessionRepo session.Repository,
	passHasher password.Hasher,
	tokenService token.Service,
) *LoginUserHandler {
	return &LoginUserHandler{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		passHasher:   passHasher,
		tokenService: tokenService,
	}
}

func (h *LoginUserHandler) Handle(ctx context.Context, cmd LoginUser) (LoginUserResponse, error) {
	u, err := h.userRepo.GetByLogin(ctx, cmd.Login)
	if err != nil && !errors.Is(err, apperr.ErrNotFound) {
		return LoginUserResponse{}, err
	}
	if errors.Is(err, apperr.ErrNotFound) {
		return LoginUserResponse{}, app.ErrInvalidCredantials
	}

	eq, err := h.passHasher.Compare(cmd.Password, u.PasswordHash)
	if err != nil {
		return LoginUserResponse{}, err
	}
	if !eq {
		return LoginUserResponse{}, app.ErrInvalidCredantials
	}

	refreshToken := h.tokenService.GenerateRefreshToken()
	refreshTokenHash := h.tokenService.HashRefreshToken(refreshToken)

	session := session.Create(
		u.ID(),
		refreshTokenHash,
		cmd.UserAgent,
		cmd.IP,
	)
	if err = h.sessionRepo.Save(ctx, session); err != nil {
		return LoginUserResponse{}, err
	}

	accessToken, err := h.tokenService.GenerateAccessToken(session)
	if err != nil {
		return LoginUserResponse{}, err
	}

	return LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
