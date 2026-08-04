package domain

import (
	"errors"
)

var ErrAlreadyExists = errors.New("already exists")
var ErrDeleted = errors.New("deleted")
var ErrNotFound = errors.New("not found")
