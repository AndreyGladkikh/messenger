package query

import (
	"context"

	"uuid"
)

const GetAccountQueryName = "GetAccount"

type GetAccount struct {
	ID uuid.UUID
}

func (q GetAccount) IsQuery() {}

func (q GetAccount) Name() string {
	return GetAccountQueryName
}

type GetAccountHandler struct {
	reader Reader
}

func NewGetAccountHandler(
	reader Reader,
) *GetAccountHandler {
	return &GetAccountHandler{
		reader: reader,
	}
}

func (h *GetAccountHandler) Handle(ctx context.Context, q GetAccount) (AccountView, error) {
	return h.reader.GetAccount(ctx, q.ID)
}
