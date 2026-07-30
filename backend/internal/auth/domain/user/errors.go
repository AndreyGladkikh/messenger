package user

import "errors"

var (
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrWrongPassword = errors.New("wrong password")
)
