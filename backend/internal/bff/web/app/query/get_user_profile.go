package query

import (
	"context"
	"messenger/messenger/internal/bff/web/app/port"
	"messenger/messenger/internal/bff/web/app/view"
	"messenger/messenger/internal/shared/app/apperr"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

const GetUserProfileQueryName = "GetUserProfile"

var _ query.Query = GetUserProfile{}

type GetUserProfile struct {
	UserID uuid.UUID
}

func (g GetUserProfile) Name() string {
	return GetUserProfileQueryName
}

type GetUserProfileHandler struct {
	directoryReader port.DirectoryReader
	mediaReader     port.MediaReader
}

func NewGetUserProfileHandler(
	directoryReader port.DirectoryReader,
	mediaReader port.MediaReader,
) *GetUserProfileHandler {
	return &GetUserProfileHandler{
		directoryReader: directoryReader,
		mediaReader:     mediaReader,
	}
}

func (h *GetUserProfileHandler) Handle(ctx context.Context, qry GetUserProfile) (view.UserProfile, error) {
	profiles, err := h.directoryReader.GetUserProfiles(ctx, []uuid.UUID{qry.UserID})
	if err != nil {
		return view.UserProfile{}, err
	}

	if len(profiles) == 0 {
		return view.UserProfile{}, apperr.NewNotFound("Profile not found")
	}

	p := profiles[0]

	var avatarURL *string
	if p.AvatarID != nil {
		urls, err := h.mediaReader.GetAvatarURLs(ctx, []uuid.UUID{*p.AvatarID})
		if err != nil {
			return view.UserProfile{}, err
		}

		if u, ok := urls[*p.AvatarID]; ok {
			avatarURL = &u
		}
	}

	return view.UserProfile{
		UserID:      p.UserID,
		Handle:      p.Handle,
		DisplayName: p.DisplayName,
		Bio:         p.Bio,
		Birthday:    p.Birthday,
		AvatarURL:   avatarURL,
	}, nil
}
