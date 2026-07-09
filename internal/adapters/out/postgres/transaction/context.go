package transaction

import (
	"context"
	"database/sql"
)

type key int

const txKey key = 0

func NewContext(ctx context.Context, tx *sql.Tx) context.Context{
	return context.WithValue(ctx, txKey, tx)
}

func FromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}