package command

import (
	"context"
	"errors"
	"messenger/messenger/internal/identity/app"
	"messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/platform/repository"
)

const RefreshUserTokensCommandName = "RefreshUserTokens"

type RefreshUserTokens struct {
	RefreshToken string
}

func (c RefreshUserTokens) IsCommand() {}

func (c RefreshUserTokens) Name() string {
	return RefreshUserTokensCommandName
}

type RefreshUserTokensResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshUserTokensHandler struct {
	sessionRepo  session.Repository
	tokenService token.Service
}

func NewRefreshUserTokensHandler(
	sessionRepo session.Repository,
	tokenService token.Service,
) *RefreshUserTokensHandler {
	return &RefreshUserTokensHandler{
		sessionRepo:  sessionRepo,
		tokenService: tokenService,
	}
}

func (h *RefreshUserTokensHandler) Handle(ctx context.Context, cmd RefreshUserTokens) (RefreshUserTokensResponse, error) {
	ses, err := h.sessionRepo.GetByRefreshTokenHash(ctx, h.tokenService.HashRefreshToken(cmd.RefreshToken))
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return RefreshUserTokensResponse{}, err
	}
	if errors.Is(err, repository.ErrNotFound) {
		return RefreshUserTokensResponse{}, app.ErrInvalidCredantials
	}

	refreshToken := h.tokenService.GenerateRefreshToken()
	refreshTokenHash := h.tokenService.HashRefreshToken(refreshToken)

	ses.Refresh(refreshTokenHash)

	if err = h.sessionRepo.Save(ctx, ses); err != nil {
		return RefreshUserTokensResponse{}, err
	}

	accessToken, err := h.tokenService.GenerateAccessToken(ses)
	if err != nil {
		return RefreshUserTokensResponse{}, err
	}

	return RefreshUserTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
