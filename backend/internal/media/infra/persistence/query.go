package persistence

import (
	"context"
	"fmt"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/platform/store"
	"uuid"
)

var _ port.MediaReader = (*QueryService)(nil)

type QueryService struct {
	*store.Store
}

func NewQueryService(
	store *store.Store,
) *QueryService {
	return &QueryService{
		Store: store,
	}
}

func (s *QueryService) MediaAssetExists(ctx context.Context, mediaID uuid.UUID) (bool, error) {
	exists, err := s.Queries(ctx).MediaAssetExists(ctx, mediaID)
	if err != nil {
		return false, fmt.Errorf("media exists: %w", err)
	}

	return exists, nil
}
