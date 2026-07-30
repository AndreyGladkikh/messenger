package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToDBUUID(v string) pgtype.UUID {
	var pgUUID pgtype.UUID

	if v == "" {
		return pgUUID
	}

	err := pgUUID.Scan(v)
	if err != nil {
		panic(fmt.Errorf("failed to convert string to pgtype.UUID: %w", err))
	}

	return pgUUID
}

func ToDBText(v string) pgtype.Text {
	var pgText pgtype.Text

	if v == "" {
		return pgText
	}

	err := pgText.Scan(v)
	if err != nil {
		panic(fmt.Errorf("failed to convert string to pgtype.Text: %w", err))
	}

	return pgText
}

// func ToDBOptionalUUID(value string) uuid.NullUUID {
// 	if value == "" {
// 		return uuid.NullUUID{}
// 	}

// 	v := uuid.MustParse(value)

// 	return uuid.NullUUID{
// 		UUID:  v,
// 		Valid: true,
// 	}
// }

func ToDBOptionalUUID(id uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{
		UUID: id,
		Valid: id != uuid.Nil,
	}
}

func FromDBOptionalUUID(id uuid.NullUUID) uuid.UUID {
	return id.UUID
}

func ToDBOptionalString(v string) sql.NullString {
	if v == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: v,
		Valid:  true,
	}
}

func ToDBOptionalTime(v time.Time) sql.NullTime {
	if v.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  v,
		Valid: true,
	}
}

func FromDBUUID(value uuid.NullUUID) string {
	if value.Valid {
		return value.UUID.String()
	}
	return ""
}

func ToDBTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: !t.IsZero(),
	}
}