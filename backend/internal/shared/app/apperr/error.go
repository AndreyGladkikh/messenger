package apperr

import (
	"errors"
	"maps"
)

var ErrNotFound = errors.New("not found")

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
	Causes  []error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return errors.Join(e.Causes...)
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
	opts = append(opts, WithKind(KindNotFound), WithCause(ErrNotFound))
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

func WithCause(causes ...error) Option {
	return func(err *Error) {
		err.Causes = append(err.Causes, causes...)
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
