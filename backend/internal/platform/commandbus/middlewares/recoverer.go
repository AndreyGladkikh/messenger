package middlewares

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/commandbus"
	"runtime/debug"
)

type PanicError struct {
	Value any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}

func Recoverer(next commandbus.CommandHandler) commandbus.CommandHandler {
	fn := func(ctx context.Context, command commandbus.Command) (response any, err error) {
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
	}

	return commandbus.HandlerFunc(fn)
}
