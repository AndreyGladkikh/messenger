package query

import (
	"context"
	"uuid"
)

type Search struct {
	Str string
}

type SearchResponse struct {
	Global   []GlobalSearchResultPeerItem
	Messages []PeerItem
}

type GlobalSearchResultPeerItem struct {
	ID           uuid.UUID
	Type         string
	AvatarUrl string
	Title        string
	Handle string
}

type SearchHandler struct {
}

func (s *SearchHandler) Handle(ctx context.Context, qry Search) (SearchResponse, error) {

	return SearchResponse{}, nil
}
