package middlewares

import (
	"context"
	"errors"
	"messenger/messenger/internal/adapters/out/logger"
	"messenger/messenger/internal/platform/command_bus"
)

type LoggerMiddlewareContainer struct {
	logger *logger.Logger
}

func NewLoggerMiddlewareContainer(
	logger *logger.Logger,
) *LoggerMiddlewareContainer {
	return &LoggerMiddlewareContainer{
		logger: logger,
	}
}

func (c *LoggerMiddlewareContainer) Middleware(next command_bus.Handler) command_bus.Handler {
	fn := func(ctx context.Context, command command_bus.Command) (any, error) {
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

			return nil, err
		}
		return response, nil
	}

	return command_bus.HandlerFunc(fn)
}