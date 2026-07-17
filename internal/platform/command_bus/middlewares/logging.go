package middlewares

import (
	"context"
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
	return command_bus.HandlerFunc(func(ctx context.Context, command command_bus.Command) (any, error) {
		response, err := next.Handle(ctx, command)
		if err != nil {
			c.logger.ErrorContext(
				ctx, "error occurred while handling command",
				 "command", command.Name(),
				  "error", err,
			)
			return nil, err
		}
		return response, nil
	})
}