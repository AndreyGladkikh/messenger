package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/transport/http/contenttype"
	"messenger/messenger/internal/platform/transport/http/headers"
	"messenger/messenger/internal/shared/app/apperr"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
)

type ProblemTypesRegistry struct {
	ValidationError,
	AlreadyExists,
	MissingRequestParameter,
	MissingRequestCookie,
	Blank string
}

func NewProblemTypesRegistry(cfg config.Config) ProblemTypesRegistry {
	problemType := func(typ string) string {
		return fmt.Sprintf("%s/%s", cfg.ProblemsRegistryUrl, typ)
	}

	return ProblemTypesRegistry{
		ValidationError:         problemType("validation-error"),
		AlreadyExists:           problemType("already-exists"),
		MissingRequestParameter: problemType("missing-request-parameter"),
		MissingRequestCookie:    problemType("missing-request-cookie"),
		Blank:                   "about:blank",
	}
}

const ProblemTitleServerError = "Server Error"

type ErrorHandler struct {
	cfg                  config.Config
	logger               *logger.Logger
	problemTypesRegistry ProblemTypesRegistry
}

func NewErrorHandler(
	cfg config.Config,
	logger *logger.Logger,
	problemTypesRegistry ProblemTypesRegistry,
) *ErrorHandler {
	return &ErrorHandler{
		cfg:                  cfg,
		logger:               logger,
		problemTypesRegistry: problemTypesRegistry,
	}
}

func (h *ErrorHandler) Handle(w http.ResponseWriter, r *http.Request, err error) {
	var pd ProblemDetails
	if convert := h.converter(err); convert != nil {
		pd = convert(err)
	} else {
		pd = h.DefaultProblemDetails()
	}

	w.Header().Set(headers.ContentType, contenttype.ApplicationProblemPlusJson)
	w.WriteHeader(*pd.Status)
	e := json.NewEncoder(w).Encode(pd)
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

func (h *ErrorHandler) DefaultProblemDetails() ProblemDetails {
	statusInternalServerError := http.StatusInternalServerError

	return ProblemDetails{
		Type:   &h.problemTypesRegistry.Blank,
		Status: &statusInternalServerError,
		Title:  http.StatusText(http.StatusInternalServerError),
	}
}

func (h *ErrorHandler) converter(err error) func(err error) ProblemDetails {
	if _, ok := errors.AsType[*apperr.Error](err); ok {
		return h.createFromAppError
	} else if _, ok := errors.AsType[*openapi3filter.RequestError](err); ok {
		return h.createFromRequestError
	}

	return nil
}

func (h *ErrorHandler) createFromAppError(err error) ProblemDetails {
	if err, ok := errors.AsType[*apperr.Error](err); ok {
		status := statusByAppErr(err)
		typ, title := h.typeTitleByAppErr(err)
		problemErrors := problemErrors(err)
		var detail, code string
		if err, ok := errors.AsType[*apperr.Error](err); ok {
			detail, code = err.Message, err.Code
		}

		return ProblemDetails{
			Type:   &typ,
			Title:  title,
			Status: &status,
			Detail: &detail,
			Code:   &code,
			Errors: &problemErrors,
		}
	}
	return ProblemDetails{}
}

func (h *ErrorHandler) createFromRequestError(err error) ProblemDetails {
	if err, ok := errors.AsType[*openapi3filter.RequestError](err); ok {
		status := http.StatusBadRequest
		typ := h.problemTypesRegistry.ValidationError
		title := "Validation error"
		detail := err.Reason

		return ProblemDetails{
			Type:   &typ,
			Title:  title,
			Status: &status,
			Detail: &detail,
			// Code: &code,
			// Errors: &problemErrors,
		}
	}
	return ProblemDetails{}
}

// func (h *ErrorHandler) problemTypeUrl(typ string) string {
// 	return fmt.Sprintf("%s/problems-registry/%s", h.cfg.FrontendUrl, typ)
// }

func problemType(typ string) func(problemsRegistryUrl string) string {
	return func(problemsRegistryUrl string) string {
		return fmt.Sprintf("%s/%s", problemsRegistryUrl, typ)
	}
}

func statusByAppErr(err *apperr.Error) int {
	switch err.Kind {
	case apperr.KindInvalidArgument:
		return http.StatusBadRequest
	case apperr.KindUnauthenticated:
		return http.StatusUnauthorized
	case apperr.KindPermissionDenied:
		return http.StatusForbidden
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindAlreadyExists:
		return http.StatusConflict
	case apperr.KindDomain:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// func problemTitle(typ string) string {
// 	switch typ {
// 	case problemTypeNotFound:
// 		return "Not Found"
// 	case problemTypeValidationError:
// 		return "Request is not valid"
// 	case problemTypeServerError:
// 		return "Server error"
// 	default:
// 		return typ
// 	}
// }

func (h *ErrorHandler) typeTitleByAppErr(err error) (typ, title string) {
	if err, ok := errors.AsType[*apperr.Error](err); ok {
		switch err.Kind {
		case apperr.KindUnauthenticated:
			return h.problemTypesRegistry.Blank, err.Kind.String()
		case apperr.KindPermissionDenied:
			return h.problemTypesRegistry.Blank, err.Kind.String()
		case apperr.KindNotFound:
			return h.problemTypesRegistry.Blank, err.Kind.String()
		case apperr.KindAlreadyExists:
			return h.problemTypesRegistry.AlreadyExists, err.Kind.String()
		}
	}

	return h.problemTypesRegistry.Blank, "Server Error"
}

func problemErrors(err error) []ErrorDetail {
	var errs []ErrorDetail

	if err, ok := errors.AsType[*apperr.Error](err); ok {
		switch err.Kind {
		case apperr.KindInvalidArgument:
			errs = append(errs, ErrorDetail{
				// Detail: ,
			})
		}
	}

	return errs
}
