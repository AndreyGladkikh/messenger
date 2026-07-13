package http_server

import (
	"errors"
	"messenger/messenger/internal/domain/derr"
	"net/http"
)

func getErrorStatus(err error) int {
	if errors.Is(err, derr.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
