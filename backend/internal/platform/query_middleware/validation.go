package query_middleware

import (
	"context"
	appQuery "messenger/messenger/internal/shared/app/query"
)

type Validator interface {
	Validate(context.Context, any) error
}

type Validation struct {
	// validator Validator
}

func NewValidation(
// validator Validator,
) *Validation {
	return &Validation{
		// validator: validator,
	}
}

func (v *Validation) Wrap[Q appQuery.Query, R any](next appQuery.Handler[Q, R]) appQuery.Handler[Q, R] {
	return appQuery.HandlerFunc[Q, R](func(ctx context.Context, qry Q) (R, error) {
		// if err := v.validator.Validate(ctx, qry); err != nil {
		// 	var zero R
		// 	return zero, err
		// }

		return next.Handle(ctx, qry)
	})
}
