package http_server

import (
	"encoding/json"
	"messenger/messenger/internal/platform/apperr"
	"net/http"
)

const defaultContentType = "application/json"
const defaultHttpStatus = http.StatusOK

type Response struct {
	response    any
	err         error
	status      string
	contentType string
	httpStatus  int
}

func NewResponse(response any, err error, opts ...Option) *Response {
	httpStatus := defaultHttpStatus
	status := "success"
	if err != nil {
		err = apperr.Translate(err)

		httpStatus = getErrorStatus(err)
		status = "error"
	}

	r := &Response{
		response:    response,
		err:         err,
		status:      status,
		contentType: defaultContentType,
		httpStatus:  httpStatus,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Response) WriteTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", r.contentType)
	w.WriteHeader(r.httpStatus)

	response := map[string]any{
		"status": r.status,
	}
	if r.err == nil {
		response["data"] = r.response
	} else {
		response["error"] = r.err
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type Option func(*Response)

func WithHttpStatus(status int) Option {
	return func(r *Response) {
		r.httpStatus = status
	}
}

func WithContentType(contentType string) Option {
	return func(r *Response) {
		r.contentType = contentType
	}
}
