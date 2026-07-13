package aperr

import (
	"errors"
	"messenger/messenger/internal/domain/derr"
)

type Error struct {
	Code    string
	Message string
	Details []any
	Wrapped error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Wrapped
}

func Translate(err error) *Error {
	e := &Error{
		Wrapped: err,
	}

	if errors.Is(err, derr.ErrNotFound) {
		e.Code = "NOT_FOUND"
		e.Message = "Ресурс не найден"
	}

	return e
}
