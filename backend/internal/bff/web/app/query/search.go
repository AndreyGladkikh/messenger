package query

import (
	"context"
	"messenger/messenger/internal/bff/web/app/port"
	"messenger/messenger/internal/bff/web/app/view"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

const SearchQueryName = "web_bff.Search.v1"

var _ query.Query = Search{}

type Search struct {
	Str string
}

func (s Search) Name() string {
	return SearchQueryName
}

type SearchResponse struct {
	Contacts []view.PeerItem
	Global   []view.PeerItem
	Messages []view.PeerItem
}

type GlobalSearchResultPeerItem struct {
	ID        uuid.UUID
	Type      string
	AvatarID *uuid.UUID
	Title     string
	Handle    string
}

type SearchHandler struct {
	directoryReader port.DirectoryReader
	mediaReader port.MediaReader
}

func NewSearchHandler(
	directoryReader port.DirectoryReader,
	mediaReader port.MediaReader,
) *SearchHandler {
	return &SearchHandler{
		directoryReader: directoryReader,
		mediaReader: mediaReader,
	}
}

func (h *SearchHandler) Handle(ctx context.Context, qry Search) (SearchResponse, error) {
	results, err := h.directoryReader.Search(ctx, qry.Str)
	if err != nil {
		return SearchResponse{}, err
	}

	var avatarIDs []uuid.UUID
	for _, r := range results {
		if r.AvatarID != nil {
			avatarIDs = append(avatarIDs, *r.AvatarID)
		}
	}

	avatarURLsByID, err := h.mediaReader.GetAvatarURLs(ctx, avatarIDs)

	peerItems := make([]view.PeerItem, 0, len(results))
	for _, r := range results {
		var avatarURL *string
		if url, ok := avatarURLsByID[*r.AvatarID]; ok {
			avatarURL = &url
		}
		
		peerItems = append(peerItems, view.PeerItem{
			ID: r.ID,
			Type: r.Type,
			Title: r.DisplayName,
			Handle: r.Handle,
			AvatarURL: avatarURL,
		})
	}

	return SearchResponse{
		Global: peerItems,
	}, nil
}
