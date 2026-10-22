package command

import (
	"context"
	"fmt"
	"messenger/messenger/internal/identity/app"
	"messenger/messenger/internal/identity/app/password"
	"messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/identity/domain/model/user"
	"net/netip"
)

const RegisterUserCommandName = "RegisterUser"

type RegisterUser struct {
	Login     string
	Password  string
	UserAgent string
	IP        netip.Addr
}

func (c RegisterUser) IsCommand() {}

func (c RegisterUser) Name() string {
	return RegisterUserCommandName
}

type RegisterUserResponse struct {
	AccessToken  string
	RefreshToken string
}

type RegisterUserHandler struct {
	userRepo     user.Repository
	sessionRepo  session.Repository
	passHasher   password.Hasher
	tokenService token.Service
}

func NewRegisterUserHandler(
	userRepo user.Repository,
	sessionRepo session.Repository,
	passHasher password.Hasher,
	tokenService token.Service,
) *RegisterUserHandler {
	return &RegisterUserHandler{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		passHasher:   passHasher,
		tokenService: tokenService,
	}
}

func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUser) (RegisterUserResponse, error) {
	exists, err := h.userRepo.ExistsByLogin(ctx, cmd.Login)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	if exists {
		return RegisterUserResponse{}, app.ErrUserAlreadyExists
	}

	passHash, err := h.passHasher.Hash(cmd.Password)
	if err != nil {
		return RegisterUserResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user := user.New(
		cmd.Login,
		passHash,
	)

	err = h.userRepo.Save(ctx, user)
	if err != nil {
		return RegisterUserResponse{}, err
	}

	refreshToken := h.tokenService.GenerateRefreshToken()
	refreshTokenHash := h.tokenService.HashRefreshToken(refreshToken)

	ses := session.Create(
		user.ID(),
		refreshTokenHash,
		cmd.UserAgent,
		cmd.IP,
	)
	if err = h.sessionRepo.Save(ctx, ses); err != nil {
		return RegisterUserResponse{}, err
	}

	accessToken, err := h.tokenService.GenerateAccessToken(ses)
	if err != nil {
		return RegisterUserResponse{}, err
	}

	return RegisterUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
