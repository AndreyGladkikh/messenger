package query

import "context"

type Search struct {
	Name string
}

type SearchResponse struct {
}

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, s Search) (SearchResponse, error) {
	return SearchResponse{}, nil
}
