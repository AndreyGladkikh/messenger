package http_server

import (
	"errors"
	"messenger/messenger/internal/auth/domain/user"
	"messenger/messenger/internal/auth/infrastructure/token"
	sharedDomain "messenger/messenger/internal/shared/domain"
	"net/http"
)

func getErrorStatus(err error) int {
	switch {
	case errors.Is(err, token.ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, user.ErrLoginAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, user.ErrWrongPassword):
		return http.StatusUnauthorized
	case errors.Is(err, sharedDomain.ErrDeleted):
		return http.StatusGone
	case errors.Is(err, sharedDomain.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, sharedDomain.ErrNotFound):
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
