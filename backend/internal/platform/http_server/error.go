package http_server

import (
	"errors"
	"messenger/messenger/internal/auth/domain/user"
	sharedDomain "messenger/messenger/internal/messaging/domain"
	"messenger/messenger/internal/messaging/domain/chat"
	"net/http"
)

func getErrorStatus(err error) int {
	switch {
	case errors.Is(err, user.ErrLoginAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, user.ErrWrongPassword):
		return http.StatusUnauthorized
	case errors.Is(err, chat.ErrDeleted):
		return http.StatusGone
	case errors.Is(err, chat.ErrPrivateChatAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, sharedDomain.ErrNotFound):
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
