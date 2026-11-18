package query

import (
	"context"
	"maps"
	"messenger/messenger/internal/bff/web/app/port"
	"messenger/messenger/internal/bff/web/app/view"
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
	Peers []view.PeerItem
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
	var directParticipantEntryRefs []uuid.UUID
	for _, c := range cnvs {
		cnvsByID[c.ID] = c
		if c.AvatarID != nil {
			cnvIDsByAvatarID[*c.AvatarID] = c.ID
		}

		if c.Type == port.ConversationTypeDirect {
			cnvIDsByDirectParticipantID[*c.DirectParticipantID] = c.ID
			directParticipantEntryRefs = append(directParticipantEntryRefs, *c.DirectParticipantID)
		}
	}

	if len(directParticipantEntryRefs) > 0 {
		dirParticipants, err := h.directoryReader.GetUserProfiles(ctx, directParticipantEntryRefs)
		if err != nil {
			return HomePage{}, err
		}

		for _, dp := range dirParticipants {
			cnvID := cnvIDsByDirectParticipantID[dp.UserID]
			cnv := cnvsByID[cnvID]
			cnv.Title = &dp.DisplayName
			cnv.AvatarID = dp.AvatarID
			cnvsByID[cnvID] = cnv

			if dp.AvatarID != nil {
				cnvIDsByAvatarID[*dp.AvatarID] = cnvIDsByDirectParticipantID[dp.UserID]
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

	peers := make([]view.PeerItem, 0, len(cnvs))
	for _, c := range cnvsByID {
		var avatarURL *string
		if c.AvatarID != nil {
			if u, ok := avatarURLsByID[*c.AvatarID]; ok {
				avatarURL = &u
			}
		}

		peer := view.PeerItem{
			ID:          c.ID,
			Type:        port.ConversationTypeDirect,
			AvatarURL:   avatarURL,
			Title:       *c.Title,
			Preview:     &c.LastMessage.Body,
			PreviewTime: &c.LastMessage.SentAt,
		}

		peers = append(peers, peer)
	}

	return HomePage{
		Peers: peers,
	}, nil
}
