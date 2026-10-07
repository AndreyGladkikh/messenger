package user

import (
	"errors"
	"fmt"
	"messenger/messenger/internal/shared/domain"
)

var (
	ErrNotFound           = fmt.Errorf("user: %w", domain.ErrNotFound)
	ErrLoginAlreadyExists = fmt.Errorf("login: %w", domain.ErrAlreadyExists)
	ErrWrongPassword      = errors.New("wrong password")
)
