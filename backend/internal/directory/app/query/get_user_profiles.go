package query

import (
	"context"
	"messenger/messenger/internal/directory/app/port"
	"messenger/messenger/internal/directory/app/view"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

const GetUserProfilesQueryName = "directory.GetUserProfiles.v1"

var _ query.Query = GetUserProfiles{}

type GetUserProfiles struct {
	UserIDs []uuid.UUID
}

func (g GetUserProfiles) Name() string {
	return GetUserProfilesQueryName
}

type GetUserProfilesHandler struct {
	directoryReader port.DirectoryReader
}

func NewGetUserProfilesHandler(
	directoryReader port.DirectoryReader,
) *GetUserProfilesHandler {
	return &GetUserProfilesHandler{
		directoryReader: directoryReader,
	}
}

func (h *GetUserProfilesHandler) Handle(ctx context.Context, qry GetUserProfiles) ([]view.UserProfile, error) {
	entries, err := h.directoryReader.GetUserProfilesByIDs(ctx, qry.UserIDs)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
