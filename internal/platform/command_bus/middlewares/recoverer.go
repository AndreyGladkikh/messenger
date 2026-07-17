package middlewares

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/command_bus"
	"runtime/debug"
)

// func NewRecovererContainer(
// 	logger *logger.Logger,
// ) *LoggerMiddlewareContainer {
// 	return &LoggerMiddlewareContainer{
// 		logger: logger,
// 	}
// }

type PanicError struct {
	Value any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}

func Recoverer(next command_bus.Handler) command_bus.Handler {
	fn := func(ctx context.Context, command command_bus.Command) (response any, err error) {
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

	return command_bus.HandlerFunc(fn)
}