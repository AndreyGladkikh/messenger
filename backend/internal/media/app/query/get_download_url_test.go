package query

import (
	"context"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/media/domain/model/media"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDownloadURLUsesStoredMediaName(t *testing.T) {
	asset := media.Rehydrate(uuid.New(), "Фото с отпуска.jpg", media.StatusReady)
	ref := port.MediaLocationRef{Bucket: "media", Key: uuid.New().String()}
	storage := &downloadStorageStub{}
	handler := NewGetDownloadURLHandler(
		&downloadMediaRepositoryStub{asset: asset},
		&downloadLocationStoreStub{ref: ref},
		storage,
	)

	downloadURL, err := handler.Handle(t.Context(), GetDownloadURL{MediaID: asset.ID()})

	require.NoError(t, err)
	assert.Equal(t, "https://storage.example.com/download", downloadURL)
	assert.Equal(t, ref, storage.ref)
	assert.Equal(t, asset.Name, storage.fileName)
}

type downloadMediaRepositoryStub struct {
	media.Repository
	asset *media.Media
}

func (r *downloadMediaRepositoryStub) Get(context.Context, uuid.UUID) (*media.Media, error) {
	return r.asset, nil
}

type downloadLocationStoreStub struct {
	port.MediaLocationStore
	ref port.MediaLocationRef
}

func (s *downloadLocationStoreStub) Get(context.Context, uuid.UUID) (port.MediaLocationRef, error) {
	return s.ref, nil
}

type downloadStorageStub struct {
	port.MediaStorage
	ref      port.MediaLocationRef
	fileName string
}

func (s *downloadStorageStub) PresignDownload(_ context.Context, ref port.MediaLocationRef, fileName string) (port.DownloadInfo, error) {
	s.ref = ref
	s.fileName = fileName
	return port.DownloadInfo{PresignedURL: "https://storage.example.com/download"}, nil
}
