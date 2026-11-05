package querymiddleware

import (
	"context"
	"messenger/messenger/internal/platform/logger"
	appQuery "messenger/messenger/internal/shared/app/query"
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

func (l *Logging) Wrap[Q appQuery.Query, R any](next appQuery.Handler[Q, R]) appQuery.Handler[Q, R] {
	return appQuery.HandlerFunc[Q, R](func(ctx context.Context, qry Q) (R, error) {
		result, err := next.Handle(ctx, qry)

		logArgs := []any{
			"query", qry.Name(),
		}

		if err != nil {
			l.logger.ErrorContext(ctx, "query handling failed", append(logArgs, "error", err)...)
			return result, err
		}

		l.logger.InfoContext(ctx, "query handled", logArgs...)

		return result, nil
	})
}
