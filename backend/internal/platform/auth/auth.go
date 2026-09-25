package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

const (
	SubjectTypeUser SubjectType = "user"
	SubjectTypeBot  SubjectType = "bot"

	MethodBearerToken Method = "bearerHttpAuthentication"
	MethodApiKey  Method = "apiKey"
)

var ErrNoAuthInContext = errors.New("no auth in context")

type Auth struct {
	Method    Method
	Principal Principal
	SessionID uuid.UUID
}

type Method string

type Principal struct {
	Subject Subject
	Scope   []string
}

type Subject struct {
	ID   uuid.UUID
	Type SubjectType
}

type SubjectType string

type key int

// const principalKey key = 0

// func NewContextWithPrincipal(ctx context.Context, p Principal) context.Context {
// 	return context.WithValue(ctx, principalKey, p)
// }

// func PrincipalFromContext(ctx context.Context) (Principal, error) {
// 	p, ok := ctx.Value(principalKey).(Principal)
// 	if !ok {
// 		return Principal{}, errors.New("no principal in context")
// 	}
// 	return p, nil
// }

const authKey key = 0

func NewContextWithAuth(ctx context.Context, a Auth) context.Context {
	return context.WithValue(ctx, authKey, a)
}

func AuthFromContext(ctx context.Context) (Auth, error) {
	a, ok := ctx.Value(authKey).(Auth)
	if !ok {
		return Auth{}, ErrNoAuthInContext
	}
	return a, nil
}
