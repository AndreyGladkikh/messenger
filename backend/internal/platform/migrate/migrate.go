package migrate

import (
	"fmt"
	"messenger/messenger/internal/platform/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func New(cfg config.Config) (*migrate.Migrate, func(), error) {
	m, err := migrate.New(fmt.Sprintf("file://%s", cfg.Migrations.Path), cfg.Migrations.DBDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create migrate: %w", err)
	}

	return m, func() {
		m.Close()
	}, nil
}
