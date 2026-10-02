package reader

import (
	"context"
	"messenger/messenger/internal/bff/web/application/model/homepage"

	"uuid"
)

type Reader interface {
	GetHomePage(ctx context.Context, userID uuid.UUID) (homepage.HomePage, error)
}
