package command_middleware

import "messenger/messenger/internal/shared/app/command"

type Middleware[C command.Command, R any] interface {
	Wrap(next command.Handler[C, R]) command.Handler[C, R]
}

type Registry struct {
	Logging     *Logging
	Recovering  *Recovering
	Transaction *Transaction
}

func Wrap[C command.Command, R any](h command.Handler[C, R], r *Registry) command.Handler[C, R] {
	h = r.Logging.Wrap(h)
	h = r.Recovering.Wrap(h)
	h = r.Transaction.Wrap(h)

	return h
}
