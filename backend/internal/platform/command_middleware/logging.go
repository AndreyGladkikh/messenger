package command_middleware

import (
	"context"
	"errors"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/shared/app/command"
)

type Logging struct {
	logger *logger.Logger
}

func NewLogging(
	logger *logger.Logger,
) *Logging {
	return &Logging{
		logger: logger,
	}
}

func (c *Logging) Wrap[C command.Command, R any](next command.Handler[C, R]) command.Handler[C ,R] {
	return command.HandlerFunc[C, R](func(ctx context.Context, command C) (R, error) {
		response, err := next.Handle(ctx, command)
		if err != nil {
			logArgs := []any{
				"command", command.Name(),
				"error", err.Error(),
			}
			if err, ok := errors.AsType[*PanicError](err); ok {
				logArgs = append(logArgs, "stack", string(err.Stack))
			}

			c.logger.ErrorContext(ctx, "error occurred while handling command", logArgs...)

			var zero R
			return zero, err
		}

		return response, nil
	})
}
