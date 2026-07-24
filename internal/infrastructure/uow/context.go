package uow

import "context"

type key int

const uowKey key = 0

func NewContext(ctx context.Context, uow *UnitOfWork) context.Context {
	return context.WithValue(ctx, uowKey, uow)
}

func FromContext(ctx context.Context) (*UnitOfWork, bool) {
	uow, ok := ctx.Value(uowKey).(*UnitOfWork)
	return uow, ok
}
