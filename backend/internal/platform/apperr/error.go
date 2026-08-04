package apperr

import (
	"errors"
	authDomain "messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/auth/domain/user"
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/messaging/domain/chat"
	sharedDomain "messenger/messenger/internal/shared/domain"
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
	case errors.Is(err, user.ErrNotFound):
		return "USER_NOT_FOUND"
	case errors.Is(err, user.ErrLoginAlreadyExists):
		return "LOGIN_ALREADY_EXISTS"
	case errors.Is(err, chat.ErrDeleted):
		return "CHAT_DELETED"
	case errors.Is(err, chat.ErrPrivateChatAlreadyExists):
		return "PRIVATE_CHAT_ALREADY_EXISTS"
	case errors.Is(err, sharedDomain.ErrNotFound):
		return "NOT_FOUND"
	default:
		return ""
	}
}

func message(err error) string {
	switch {
	case errors.Is(err, authDomain.ErrInvalidCredantials):
		return "Неверный логин или пароль"
	case errors.Is(err, token.ErrInvalidToken):
		return "Невалидный токен авторизации"
	case errors.Is(err, authDomain.ErrUnauthorized):
		return "Аутентификация не удалась"
	case errors.Is(err, user.ErrNotFound):
		return "Пользователь не найден"
	case errors.Is(err, user.ErrLoginAlreadyExists):
		return "Пользователь с таким логином уже существует"
	case errors.Is(err, chat.ErrNotFound):
		return "Чат не найден"
	case errors.Is(err, chat.ErrDeleted):
		return "Чат удален"
	case errors.Is(err, chat.ErrPrivateChatAlreadyExists):
		return "Приватный чат уже существует"
	case errors.Is(err, sharedDomain.ErrNotFound):
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
