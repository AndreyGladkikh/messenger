package query

import (
	"context"
	"maps"
	"messenger/messenger/internal/bff/web/app/port"
	"messenger/messenger/internal/shared/app/query"
	"slices"

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
	mediaReader port.MediaReader
}

func NewGetHomePageHandler(
	messagingReader port.MessagingReader,
	directoryReader port.DirectoryReader,
	mediaReader port.MediaReader,
) *GetHomePageHandler {
	return &GetHomePageHandler{
		messagingReader: messagingReader,
		directoryReader: directoryReader,
		mediaReader: mediaReader,
	}
}

func (h *GetHomePageHandler) Handle(ctx context.Context, qry GetHomePage) (HomePage, error) {
	cnvs, err := h.messagingReader.ListUserConvesations(ctx, qry.UserID)
	if err != nil {
		return HomePage{}, err
	}

	cnvsByID := make(map[uuid.UUID]port.ConversationSummary)
	cnvIDsByAvatarID := make(map[uuid.UUID]uuid.UUID)
	cnvIDsByDirectParticipantID := make(map[uuid.UUID]uuid.UUID)
	var directParticipantEntryRefs []port.EntryRef
	for _, c := range cnvs {
		cnvsByID[c.ID] = c
		if c.AvatarID != nil {
			cnvIDsByAvatarID[*c.AvatarID] = c.ID
		}

		if c.Type == port.ConversationTypeDirect {
			cnvIDsByDirectParticipantID[*c.DirectParticipantID] = c.ID
			directParticipantEntryRefs = append(directParticipantEntryRefs, port.EntryRef{
				ID:   *c.DirectParticipantID,
				Type: port.DirectoryEntryTypeUser,
			})
		}
	}

	if len(directParticipantEntryRefs) > 0 {
		dirParticipants, err := h.directoryReader.GetEntries(ctx, directParticipantEntryRefs)
		if err != nil {
			return HomePage{}, err
		}

		for _, dp := range dirParticipants {
			cnv := cnvsByID[cnvIDsByDirectParticipantID[dp.ID]]
			cnv.Title = &dp.DisplayName
			cnv.AvatarID = dp.AvatarID

			if dp.AvatarID != nil {
				cnvIDsByAvatarID[*dp.AvatarID] = cnvIDsByDirectParticipantID[dp.ID]
			}
		}
	}

	avatarIDs := slices.Collect(maps.Keys(cnvIDsByAvatarID))
	var avatarURLsByID map[uuid.UUID]string
	if len(avatarIDs) > 0 {
		avatarURLsByID, err = h.mediaReader.GetAvatarURLs(ctx, avatarIDs)
		if err != nil {
			return HomePage{}, err
		}
	}

	peers := make([]PeerItem, 0, len(cnvs))
	for _, c := range cnvsByID {
		var avatarURL *string
		if c.AvatarID != nil {
			if u, ok := avatarURLsByID[*c.AvatarID]; ok {
				avatarURL = &u
			}
		}
		
		peer := PeerItem{
			ID:          c.ID,
			Type:        port.ConversationTypeDirect,
			AvatarURL:   avatarURL,
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
