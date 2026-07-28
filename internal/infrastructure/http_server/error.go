package http_server

import (
	"errors"
	"messenger/messenger/internal/domain"
	"net/http"
)

func getErrorStatus(err error) int {
	if errors.Is(err, domain.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
