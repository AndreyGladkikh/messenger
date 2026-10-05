package user

import (
	"errors"
	"fmt"
	sharedDomain "messenger/messenger/internal/shared/domain"
)

var (
	ErrNotFound           = fmt.Errorf("user: %w", sharedDomain.ErrNotFound)
	ErrLoginAlreadyExists = fmt.Errorf("login: %w", sharedDomain.ErrAlreadyExists)
	ErrWrongPassword      = errors.New("wrong password")
)
