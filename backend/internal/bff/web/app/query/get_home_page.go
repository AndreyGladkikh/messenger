package query

import (
	"context"
	"messenger/messenger/internal/bff/web/app/port"
	"messenger/messenger/internal/shared/app/query"

	"uuid"
)

const GetHomePageQueryName = "GetHomePage"

var _ query.Query = GetHomePage{}

type GetHomePage struct {
	UserID uuid.UUID
}

func (g GetHomePage) Name() string {
	return GetHomePageQueryName
}

type HomePage struct {
	Peers []PeerItem
}

type GetHomePageHandler struct {
	messagingReader port.MessagingReader
	directoryReader port.DirectoryReader
}

func NewGetHomePageHandler(
	messagingReader port.MessagingReader,
	directoryReader port.DirectoryReader,
) *GetHomePageHandler {
	return &GetHomePageHandler{
		messagingReader: messagingReader,
		directoryReader: directoryReader,
	}
}

func (h *GetHomePageHandler) Handle(ctx context.Context, qry GetHomePage) (HomePage, error) {
	cnvs, err := h.messagingReader.ListUserConvesations(ctx, qry.UserID)
	if err != nil {
		return HomePage{}, err
	}

	cnvIdxMap := make(map[uuid.UUID]int)
	var entryRefs []port.EntryRef
	for i, c := range cnvs {
		if c.Type == port.ConversationTypeDirect {
			cnvIdxMap[*c.DirectParticipantID] = i
			entryRefs = append(entryRefs, port.EntryRef{
				ID:   *c.DirectParticipantID,
				Type: port.DirectoryEntryTypeUser,
			})
		}
	}

	if len(entryRefs) > 0 {
		entries, err := h.directoryReader.GetEntries(ctx, entryRefs)
		if err != nil {
			return HomePage{}, err
		}

		for _, e := range entries {
			cnvs[cnvIdxMap[e.ID]].Title = &e.DisplayName
			cnvs[cnvIdxMap[e.ID]].AvatarURL = e.AvatarURL
		}
	}

	peers := make([]PeerItem, 0, len(cnvs))
	for _, c := range cnvs {
		peer := PeerItem{
			ID:          c.ID,
			Type:        port.ConversationTypeDirect,
			AvatarUrl:   c.AvatarURL,
			Title:       *c.Title,
			Preview:     c.LastMessage.Body,
			PreviewTime: c.LastMessage.SentAt,
		}

		peers = append(peers, peer)
	}

	return HomePage{
		Peers: peers,
	}, nil
}
