package reader

import (
	"context"
	"messenger/messenger/internal/bff/web/application/model/homepage"

	"github.com/google/uuid"
)

type Reader interface {
	GetHomePage(ctx context.Context, userID uuid.UUID) (homepage.HomePage, error)
}
