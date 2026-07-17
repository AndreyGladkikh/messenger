package mapping

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

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

func OptionalString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid: true,
	}
}

func OptionalTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time: value,
		Valid: true,
	}
}

func UUIDString(value uuid.NullUUID) string {
	if value.Valid {
		return value.UUID.String()
	}
	return ""
}