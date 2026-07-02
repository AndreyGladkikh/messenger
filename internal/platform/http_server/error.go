package http_server

import (
	"errors"
	"messenger/messenger/internal/domain/domainerr"
	"net/http"
)

func getErrorStatus(err error) int {
	if errors.Is(err, domainerr.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}