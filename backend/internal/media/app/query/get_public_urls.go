package query

import (
	"context"
	"fmt"
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
}

func NewGetPublicURLsHandler() *GetPublicURLsHandler {
	return &GetPublicURLsHandler{}
}

func (h *GetPublicURLsHandler) Handle(ctx context.Context, qry GetPublicURLs) (map[uuid.UUID]string, error) {
	res := make(map[uuid.UUID]string)
	for _, id := range qry.MediaIDs {
		res[id] = fmt.Sprintf("%s/%s/%s/original", h.cfg.CDNHost, s3.BucketPublic, id)
	}

	return res, nil
}
