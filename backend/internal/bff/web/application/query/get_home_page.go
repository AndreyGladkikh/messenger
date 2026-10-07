package query

import (
	"context"

	"uuid"
)

type GetHomePage struct {
	UserID uuid.UUID
}

type HomePage struct {
	Peers []PeerItem
}

type GetHomePageHandler struct {
	reader Reader
}

func (h *GetHomePageHandler) Handle(ctx context.Context, qry GetHomePage) (HomePage, error) {
	return h.reader.GetHomePage(ctx, qry.UserID)
}
