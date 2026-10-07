package query

import "context"

type Search struct {
	Str string
}

type SearchResponse struct {
	// Global []GlobalSearchResult
	// Messages []Message
}

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, s Search) (SearchResponse, error) {
	return SearchResponse{}, nil
}
