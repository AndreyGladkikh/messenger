package mapping

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgUUID(v string) pgtype.UUID {
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

func PgText(v string) pgtype.Text {
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

func OptionalUUID(value string) uuid.NullUUID {
	if value == "" {
		return uuid.NullUUID{}
	}

	v := uuid.MustParse(value)

	return uuid.NullUUID{
		UUID: v,
		Valid: true,
	}
}

func OptionalString(v string) sql.NullString {
	if v == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: v,
		Valid: true,
	}
}

func OptionalTime(v time.Time) sql.NullTime {
	if v.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time: v,
		Valid: true,
	}
}

func UUIDString(value uuid.NullUUID) string {
	if value.Valid {
		return value.UUID.String()
	}
	return ""
}

func ToDBTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time: t,
		Valid: !t.IsZero(),
	}
}