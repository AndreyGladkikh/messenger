package query

import (
	"context"
	"messenger/messenger/internal/bff/web/application/model/homepage"
	"messenger/messenger/internal/bff/web/application/reader"

	"github.com/google/uuid"
)

type GetHomePage struct{
	UserID uuid.UUID
}

type GetHomePageHandler struct {
	reader reader.Reader
}

func (h *GetHomePageHandler) Handle(ctx context.Context, query GetHomePage) (homepage.HomePage, error) {
	hp, err := h.reader.GetHomePage(ctx, query.UserID)
	if err != nil {
		return homepage.HomePage{}, err
	}
	return hp, nil
}
