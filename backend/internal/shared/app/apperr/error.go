package apperr

import (
	"maps"
	"errors"
	"messenger/messenger/internal/directory/domain/entry"
)

type kind int

func (k kind) String() string {
	return kindName[k]
}

const (
	KindInvalidArgument kind = iota
	KindUnauthenticated
	KindPermissionDenied
	KindNotFound
	KindAlreadyExists
	KindDomain
)

var kindName = map[kind]string{
	KindInvalidArgument:  "Invalid argument",
	KindUnauthenticated:  "Unauthenticated",
	KindPermissionDenied: "Permission denied",
	KindNotFound:         "Not found",
	KindAlreadyExists:    "Already exists",
	KindDomain:           "Business rule violation",
}

const (
	CodeUnsupportableDirectoryEntryTargetType = "UNSUPPORTABLE_DIRECTORY_ENTRY_TARGET_TYPE"
	CodeInvalidCredentials                    = "INVALID_CREDENTIALS"
	CodeInvalidToken                          = "INVALID_TOKEN"
	CodeUnauthorized                          = "UNAUTHORIZED"
	CodeAlreadyExists                         = "ALREADY_EXISTS"
	CodeUserNotFound                          = "USER_NOT_FOUND"
	CodeLoginAlreadyExists                    = "LOGIN_ALREADY_EXISTS"
	CodeChatNotFound                          = "CHAT_NOT_FOUND"
	CodeChatDeleted                           = "CHAT_DELETED"
	CodePrivateChatAlreadyExists              = "PRIVATE_CHAT_ALREADY_EXISTS"
	CodeNotFound                              = "NOT_FOUND"
)

type Error struct {
	Kind    kind
	Message string
	Code    string
	Details map[string]any
	Wrapped error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Wrapped
}

func New(msg string, opts ...Option) *Error {
	err := &Error{
		Message: msg,
		Details: make(map[string]any),
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

func NewUnauthenticated(msg string, opts ...Option) *Error {
	opts = append(opts, WithKind(KindUnauthenticated))
	return New(msg, opts...)
}

func NewPermissionDenied(msg string, opts ...Option) *Error {
	opts = append(opts, WithKind(KindPermissionDenied))
	return New(msg, opts...)
}

func NewNotFound(msg string, opts ...Option) *Error {
	opts = append(opts, WithKind(KindNotFound))
	return New(msg, opts...)
}

func NewAlreadyExists(msg string, opts ...Option) *Error {
	opts = append(opts, WithKind(KindAlreadyExists))
	return New(msg, opts...)
}

func NewInvalidArgument(msg string, opts ...Option) *Error {
	opts = append(opts, WithKind(KindInvalidArgument))
	return New(msg, opts...)
}

type Option func(*Error)

func WithCause(cause error) Option {
	return func(err *Error) {
		err.Wrapped = cause
	}
}

func WithKind(k kind) Option {
	return func(err *Error) {
		err.Kind = k
	}
}

func WithCode(code string) Option {
	return func(err *Error) {
		err.Code = code
	}
}

func WithDetails(details map[string]any) Option {
	return func(err *Error) {
		maps.Copy(err.Details, details)
	}
}

// todo remove
func Map(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := errors.AsType[*Error](err); ok {
		return err
	}

	if _, ok := errors.AsType[entry.UnsupportableTargetTypeError](err); ok {
		return &Error{
			Code:    CodeUnsupportableDirectoryEntryTargetType,
			Message: "Unsupportable Directory Entry Target Type",
			Wrapped: err,
		}
	}

	return err
}
