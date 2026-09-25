package domain

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidCredantials = fmt.Errorf("%w: invalid credentials", ErrUnauthorized)
)
