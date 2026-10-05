package apperr

import (
	"errors"
	auth "messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/auth/domain/user"
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/shared/domain"
)

const (
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeInvalidToken = "INVALID_TOKEN"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeAlreadyExists = "ALREADY_EXISTS"
	CodeUserNotFound = "USER_NOT_FOUND"
	CodeLoginAlreadyExists = "LOGIN_ALREADY_EXISTS"
	CodeChatNotFound = "CHAT_NOT_FOUND"
	CodeChatDeleted = "CHAT_DELETED"
	CodePrivateChatAlreadyExists = "PRIVATE_CHAT_ALREADY_EXISTS"
	CodeNotFound = "NOT_FOUND"
)

type Error struct {
	Code    string
	Message string
	Details map[string]any
	Wrapped error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Wrapped
}

func Map(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := errors.AsType[*Error](err); ok {
		return err
	}

	switch {
	case errors.Is(err, token.ErrInvalidToken):
		return &Error{
			Code: CodeInvalidToken,
			Message: "Invalid access token",
			Wrapped: err,
		}
	case errors.Is(err, auth.ErrInvalidCredantials):
		return &Error{
			Code: CodeInvalidCredentials,
			Message: "Invalid credentials",
			Wrapped: err,
		}
	case errors.Is(err, auth.ErrUnauthorized):
		return &Error{
			Code: CodeUnauthorized,
			Message: "Authentication failed",
			Wrapped: err,
		}
	case errors.Is(err, user.ErrNotFound):
		return &Error{
			Code: CodeUserNotFound,
			Message: "User not found",
			Wrapped: err,
		}
	case errors.Is(err, user.ErrLoginAlreadyExists):
		return &Error{
			Code: CodeLoginAlreadyExists,
			Message: "Login already exists",
			Wrapped: err,
		}

	case errors.Is(err, domain.ErrNotFound):
		return &Error{
			Code: CodeNotFound,
			Message: "Not found",
			Wrapped: err,
		}
	case errors.Is(err, domain.ErrAlreadyExists):
		return &Error{
			Code: CodeAlreadyExists,
			Message: "Already exists",
			Wrapped: err,
		}
	default:
		return err
	}
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
	case errors.Is(err, domain.ErrNotFound):
		return "NOT_FOUND"
	default:
		return ""
	}
}

func message(err error) string {
	switch {
	case errors.Is(err, auth.ErrInvalidCredantials):
		return "Неверный логин или пароль"
	case errors.Is(err, token.ErrInvalidToken):
		return "Невалидный токен авторизации"
	case errors.Is(err, auth.ErrUnauthorized):
		return "Аутентификация не удалась"
	case errors.Is(err, user.ErrNotFound):
		return "Пользователь не найден"
	case errors.Is(err, user.ErrLoginAlreadyExists):
		return "Пользователь с таким логином уже существует"
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
