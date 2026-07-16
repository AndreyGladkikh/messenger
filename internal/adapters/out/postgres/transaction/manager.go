package transaction

import (
	"context"
	"database/sql"
)

type Manager struct {
	db *sql.DB
}

func NewManager(
	db *sql.DB,
) *Manager {
	return &Manager{
		db,
	}
}

func (t *Manager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx = NewContext(ctx, tx)

	if err = fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
