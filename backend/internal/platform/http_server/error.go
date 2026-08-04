package http_server

import (
	"errors"
	authDomain "messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/auth/domain/user"
	sharedDomain "messenger/messenger/internal/shared/domain"
	"net/http"
)

func getErrorStatus(err error) int {
	switch {
	case errors.Is(err, user.ErrLoginAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, authDomain.ErrUnauthorized):
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
