package auth

import "context"

type Service interface {
	UserExists(ctx context.Context, userID string) bool
}
