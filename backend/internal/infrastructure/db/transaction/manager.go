package transaction

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	db *pgxpool.Pool
}

func NewManager(
	db *pgxpool.Pool,
) *Manager {
	return &Manager{
		db,
	}
}

func (t *Manager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	ctx = NewContext(ctx, tx)

	if err = fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
