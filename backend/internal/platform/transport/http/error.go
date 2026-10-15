package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"messenger/messenger/internal/directory/domain/entry"
	auth "messenger/messenger/internal/identity/domain/model"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/transport/http/contenttype"
	"messenger/messenger/internal/platform/transport/http/headers"
	"messenger/messenger/internal/shared/app/apperr"
	"messenger/messenger/internal/shared/domain"
	"net/http"
)

var (
	problemTypeNotFound        = "not-found"
	problemTypeAlreadyExists   = "already-exists"
	problemTypeValidationError = "validation-error"
	problemTypeServerError     = "server-error"
)

type ErrorHandler struct {
	cfg    *config.Config
	logger *logger.Logger
}

func NewErrorHandler(
	cfg *config.Config,
	logger *logger.Logger,
) *ErrorHandler {
	return &ErrorHandler{
		cfg:    cfg,
		logger: logger,
	}
}

func (h *ErrorHandler) Handle(w http.ResponseWriter, r *http.Request, err error) {
	status := getErrorStatus(err)

	w.Header().Set(headers.ContentType, contenttype.ApplicationProblemPlusJson)
	w.WriteHeader(status)

	problemDetails := h.mapError(err)
	problemDetails.Status = &status

	e := json.NewEncoder(w).Encode(problemDetails)
	if e != nil {
		h.logger.ErrorContext(
			r.Context(),
			"failed to write error to http.ResponseWriter",
			"error", err,
			"cause", e,
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *ErrorHandler) mapError(err error) ProblemDetails {
	if err, ok := errors.AsType[*apperr.Error](err); ok {
		// var typ, typeUrl string

		// switch {
		// case errors.Is(err, sharedDomain.ErrNotFound):
		// 	typ = problemTypeNotFound
		// case errors.Is(err, sharedDomain.ErrAlreadyExists):
		// 	typ = problemTypeAlreadyExists
		// }

		// if typ != "" {
		// 	typeUrl = h.problemTypeUrl(typ)
		// }

		return ProblemDetails{
			// Type: &typeUrl,
			// Title: problemTitle(typ),
			Title: err.Error(),
			Code:  &err.Code,
		}
	}

	// typeUrl := h.problemTypeUrl(problemTypeServerError)
	return ProblemDetails{
		// Type: &typeUrl,
		// Title: problemTitle(problemTypeServerError),
		Title: "Server error",
	}
}

func (h *ErrorHandler) problemTypeUrl(typ string) string {
	return fmt.Sprintf("%s/problems-registry/%s", h.cfg.FrontendUrl, typ)
}

func getErrorStatus(err error) int {
	if _, ok := errors.AsType[entry.UnsupportableTargetTypeError](err); ok {
		return http.StatusUnprocessableEntity
	} else if (errors.Is(err, auth.ErrUnauthorized)) {
		return http.StatusUnauthorized
	} else if (errors.Is(err, auth.ErrForbidden)) {
		return http.StatusForbidden
	} else if (errors.Is(err, domain.ErrNotFound)) {
		return http.StatusNotFound
	} else if (errors.Is(err, domain.ErrDeleted)) {
		return http.StatusGone
	} else if (errors.Is(err, domain.ErrAlreadyExists)) {
		return http.StatusConflict
	}

	return http.StatusInternalServerError
}

func problemTitle(typ string) string {
	switch typ {
	case problemTypeNotFound:
		return "Not Found"
	case problemTypeValidationError:
		return "Request is not valid"
	case problemTypeServerError:
		return "Server error"
	default:
		return typ
	}
}
