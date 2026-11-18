package query

import (
	"context"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

const MediaExistsQueryName = "media.MediaExists.v1"

var _ query.Query = MediaExists{}

type MediaExists struct {
	ID uuid.UUID
}

func (q MediaExists) Name() string {
	return MediaExistsQueryName
}

type MediaExistsHandler struct {
	mediaReader port.MediaReader
}

func NewMediaExistsHandler(
	mediaReader port.MediaReader,
) *MediaExistsHandler {
	return &MediaExistsHandler{
		mediaReader: mediaReader,
	}
}

func (h *MediaExistsHandler) Handle(ctx context.Context, qry MediaExists) (bool, error) {
	exists, err := h.mediaReader.MediaAssetExists(ctx, qry.ID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
