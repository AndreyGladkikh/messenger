package query

import (
	"context"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/media/domain/model/media"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

const GetDownloadURLQueryName = "GetDownloadURL"

var _ query.Query = GetDownloadURL{}

type GetDownloadURL struct {
	MediaID uuid.UUID
}

func (g GetDownloadURL) Name() string {
	return GetDownloadURLQueryName
}

type GetDownloadURLHandler struct {
	mediaRepo          media.Repository
	mediaLocationStore port.MediaLocationStore
	objectStorage      port.MediaStorage
}

func NewGetDownloadURLHandler(
	mediaRepo media.Repository,
	mediaLocationStore port.MediaLocationStore,
	objectStorage port.MediaStorage,
) *GetDownloadURLHandler {
	return &GetDownloadURLHandler{
		mediaRepo:          mediaRepo,
		mediaLocationStore: mediaLocationStore,
		objectStorage:      objectStorage,
	}
}

func (h *GetDownloadURLHandler) Handle(ctx context.Context, qry GetDownloadURL) (string, error) {
	m, err := h.mediaRepo.Get(ctx, qry.MediaID)
	if err != nil {
		return "", err
	}

	locRef, err := h.mediaLocationStore.Get(ctx, m.ID())
	if err != nil {
		return "", err
	}

	info, err := h.objectStorage.PresignDownload(ctx, locRef, m.Name)
	if err != nil {
		return "", err
	}

	return info.PresignedURL, nil
}
