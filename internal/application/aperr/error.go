package aperr

import (
	"errors"
	"messenger/messenger/internal/domain/derr"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details []any `json:"details"`
	Wrapped error `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Wrapped
}

func Translate(err error) *Error {
	e := &Error{
		Message: "Непредвиденная ошибка",
		Details: make([]any, 0),
		Wrapped: err,
	}

	if errors.Is(err, derr.ErrNotFound) {
		e.Code = "NOT_FOUND"
		e.Message = "Ресурс не найден"
	}

	return e
}
