package repositories

import "github.com/jmoiron/sqlx"

type ChatRepository struct {
	db *sqlx.DB
}
