package command

import "errors"

var ErrValidation = errors.New("validation error")

type Command interface {
	IsCommand()
	Name() string
}