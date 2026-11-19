package query

import (
	"context"
	"fmt"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/s3"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

const GetPublicURLQueryName = "media.GetPublicURLs"

var _ query.Query = GetPublicURLs{}

type GetPublicURLs struct {
	MediaIDs []uuid.UUID
}

func (g GetPublicURLs) Name() string {
	return GetPublicURLQueryName
}

type GetPublicURLsHandler struct {
	cfg config.Config
	mediaLocationStore port.MediaLocationStore
}

func NewGetPublicURLsHandler(
	cfg config.Config,
	mediaLocationStore port.MediaLocationStore,
) *GetPublicURLsHandler {
	return &GetPublicURLsHandler{
		cfg: cfg,
		mediaLocationStore: mediaLocationStore,
	}
}

func (h *GetPublicURLsHandler) Handle(ctx context.Context, qry GetPublicURLs) (map[uuid.UUID]string, error) {
	locs, err := h.mediaLocationStore.GetMany(ctx, qry.MediaIDs)
	if err != nil {
		return nil, err
	}

	res := make(map[uuid.UUID]string, len(locs))
	for id, loc := range locs {
		if loc.Bucket == s3.BucketPublic {
			res[id] = fmt.Sprintf("%s/%s/%s", h.cfg.CDNHost, loc.Bucket, loc.Key)
		}
	}

	return res, nil
}
