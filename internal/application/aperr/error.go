package aperr

import (
	"errors"
	"messenger/messenger/internal/domain"
	"messenger/messenger/internal/domain/chat"
)

type Error struct {
	Code    string         `json:"code,omitempty"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
	Wrapped error          `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Wrapped
}

func Translate(err error) *Error {
	return &Error{
		Code:    code(err),
		Message: message(err),
		Details: details(err),
		Wrapped: err,
	}
}

func code(err error) string {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return "NOT_FOUND"
	default:
		return "INTERNAL"
	}
}

func message(err error) string {
	switch {
	case errors.Is(err, chat.ErrNotFound):
		return "Чат не найден"
	case errors.Is(err, domain.ErrNotFound):
		return "Ресурс не найден"
	default:
		return "Непредвиденная ошибка"
	}
}

func details(err error) map[string]any {
	switch {
	default:
		return nil
	}
}
