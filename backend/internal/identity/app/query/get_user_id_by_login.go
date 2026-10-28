package query

import (
	"context"
	"uuid"
)

type GetUserIDByLogin struct {
	Login string
}

type GetUserIDByLoginHandler struct {
	Reader Reader
}

func (h *GetUserIDByLoginHandler) Handle(ctx context.Context, qry GetUserIDByLogin) (uuid.UUID, error) {
	id, err := h.Reader.GetUserIDByLogin(ctx, qry.Login)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
