package domain

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidCredantials = fmt.Errorf("%w: invalid credentials", ErrUnauthorized)
)
