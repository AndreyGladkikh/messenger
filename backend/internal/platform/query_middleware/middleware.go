package query_middleware

import "messenger/messenger/internal/shared/app/query"

type Middleware[Q query.Query, R any] interface {
	Wrap(query.Handler[Q, R]) query.Handler[Q, R]
}

type Registry struct {
	Logging    *Logging
	Validation *Validation
}

func Wrap[Q query.Query, R any](handler query.Handler[Q, R], middlewares *Registry) query.Handler[Q, R] {
	handler = middlewares.Validation.Wrap(handler)
	handler = middlewares.Logging.Wrap(handler)

	return handler
}
