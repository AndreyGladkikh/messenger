package auth

import (
	"context"

	"github.com/google/uuid"
)

type key int

const (
	userIDKey key = 0
	sessionIDKey key = 1
)

func NewContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}

func NewContextWithSessionID(ctx context.Context, sessionID uuid.UUID) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

func SessionIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	sessionID, ok := ctx.Value(sessionIDKey).(uuid.UUID)
	return sessionID, ok
}
