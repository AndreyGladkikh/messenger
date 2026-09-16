package http_server

import (
	"encoding/json"
	"messenger/messenger/internal/platform/apperr"
	"net/http"
	"time"

	"github.com/coder/websocket"
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

type ResponseStatus string

const (
	ResponseStatusSuccess ResponseStatus = "success"
	ResponseStatusError ResponseStatus = "error"
)

type ResponseBody struct {
	Status ResponseStatus `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  *apperr.Error  `json:"error,omitempty"`
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

// === http

type ResponseNew struct {
	Status ResponseStatus `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  *apperr.Error  `json:"error,omitempty"`
}

func (r ResponseNew) isError() bool {
	return r.Error != nil
}

func NewResponseNew(data any, err error) ResponseNew {
	var response ResponseNew

	if err == nil {
		response.Status = ResponseStatusSuccess
		response.Data = data
	} else {
		response.Status = ResponseStatusError
		response.Error = apperr.Translate(err)
	}

	return response
}

type HTTPResponse struct {
	response ResponseNew
	contentType string
	status  int
}

func NewHTTPResponse(data any, err error, opts ...OptionNew) *HTTPResponse {
	response := NewResponseNew(data, err)

	httpResponse := &HTTPResponse{
		response: response,
		contentType: defaultContentType,
	}

	if response.isError() {
		httpResponse.status = getErrorStatus(response.Error)
	} else {
		httpResponse.status = http.StatusOK
	}

	for _, opt := range opts {
		opt(httpResponse)
	}

	return httpResponse
}

func (r *HTTPResponse) WriteTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", r.contentType)
	w.WriteHeader(r.status)

	err := json.NewEncoder(w).Encode(r.response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type OptionNew func(*HTTPResponse)

func WithStatus(status int) OptionNew {
	return func(r *HTTPResponse) {
		r.status = status
	}
}

func WithContentTypeNew(contentType string) OptionNew {
	return func(r *HTTPResponse) {
		r.contentType = contentType
	}
}

// == ws

type WSResponse struct {
	response ResponseNew
	contentType string
	timeout time.Duration
}

func NewWSResponse(data any, err error, opts ...WSResponseOption) *WSResponse {
	response := NewResponseNew(data, err)

	httpResponse := &WSResponse{
		response: response,
		contentType: defaultContentType,
	}

	for _, opt := range opts {
		opt(httpResponse)
	}

	return httpResponse
}

func (r *WSResponse) WriteTo(c *websocket.Conn) {
	// err := json.NewEncoder(w).Encode(r.response)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}

type WSResponseOption func(*WSResponse)

func WithTimeout(timeout time.Duration) WSResponseOption {
	return func(r *WSResponse) {
		r.timeout = timeout
	}
}