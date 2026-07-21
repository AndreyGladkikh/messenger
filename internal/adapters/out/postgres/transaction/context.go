package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type key int

const txKey key = 0

func NewContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func FromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	return tx, ok
}
