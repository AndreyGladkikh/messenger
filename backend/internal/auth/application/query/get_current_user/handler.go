package get_current_user

import (
	"context"

	"github.com/google/uuid"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, q *Query) (*AuthView, error) {
	return &AuthView{
		ID:   uuid.New(),
		Name: "YoungTalent",
	}, nil
}
