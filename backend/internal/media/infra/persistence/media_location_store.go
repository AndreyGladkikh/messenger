package persistence

import (
	"context"
	"fmt"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/platform/store"
	"uuid"
)

var _ port.MediaLocationStore = (*MediaLocationStore)(nil)

type MediaLocationStore struct {
	*store.Store
}

func NewMediaLocationStore(
	store *store.Store,
) *MediaLocationStore {
	return &MediaLocationStore{
		Store: store,
	}
}

func (s *MediaLocationStore) Save(ctx context.Context, mediaID uuid.UUID, ref port.MediaLocationRef) error {
	err := s.Queries(ctx).UpdateMediaAssetLocation(ctx, sqlc.UpdateMediaAssetLocationParams{
		ID:     mediaID,
		Bucket: ref.Bucket,
		Key:    ref.Key,
	})
	if err != nil {
		return fmt.Errorf("update media location: %w", err)
	}

	return nil
}

func (s *MediaLocationStore) Get(ctx context.Context, mediaID uuid.UUID) (port.MediaLocationRef, error) {
	m, err := s.Queries(ctx).GetMediaAsset(ctx, mediaID)
	if err != nil {
		return port.MediaLocationRef{}, fmt.Errorf("get media location: %w", err)
	}

	return port.MediaLocationRef{
		Bucket: m.Bucket,
		Key:    m.Key,
	}, nil
}
