package query

import (
	"context"

	"uuid"
)

type Reader interface {
	GetHomePage(ctx context.Context, userID uuid.UUID) (HomePage, error)
}
