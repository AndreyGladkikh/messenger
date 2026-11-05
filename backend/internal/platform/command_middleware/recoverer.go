package command_middleware

import (
	"context"
	"fmt"
	"messenger/messenger/internal/shared/app/command"
	"runtime/debug"
)

type PanicError struct {
	Value any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}

type Recovering struct {}

func NewRecovering() *Recovering {
	return &Recovering{}
}

func (r *Recovering) Wrap[C command.Command, R any](next command.Handler[C, R]) command.Handler[C ,R] {
	return command.HandlerFunc[C, R](func(ctx context.Context, command C) (response R, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = &PanicError{
					Value: r,
					Stack: debug.Stack(),
				}
			}
		}()

		response, err = next.Handle(ctx, command)
		return
	})
}
