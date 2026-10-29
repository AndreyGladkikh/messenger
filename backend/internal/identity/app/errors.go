package app

import "messenger/messenger/internal/shared/app/apperr"

var (
	ErrInvalidCredantials = apperr.NewUnauthenticated("Invalid credentials")
	ErrNeedLogin          = apperr.NewUnauthenticated("You need to log in")
	ErrUserAlreadyExists  = apperr.NewAlreadyExists("User already exists")
)
