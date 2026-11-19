package port

import (
	"context"
	"uuid"
)

type UploadInfo struct {
	PresignedURL string
	LocationRef  MediaLocationRef
}

type DownloadInfo struct {
	PresignedURL string
}

type MediaInfo struct {
}

type MediaLocationRef struct {
	Bucket string
	Key    string
}

type AccessLevel string

const (
	AccessLevelPrivate AccessLevel = "private"
	AccessLevelPublic  AccessLevel = "public"
)

type MediaStorage interface {
	PresignUpload(ctx context.Context, uploadID uuid.UUID, accessLevel AccessLevel) (UploadInfo, error)
	PresignDownload(ctx context.Context, ref MediaLocationRef, fileName string) (DownloadInfo, error)
	GetObjectInfo(ctx context.Context, ref MediaLocationRef) (MediaInfo, error)
}

type MediaLocationStore interface {
	Save(ctx context.Context, mediaID uuid.UUID, ref MediaLocationRef) error
	Get(ctx context.Context, mediaID uuid.UUID) (MediaLocationRef, error)
	GetMany(ctx context.Context, mediaIDs []uuid.UUID) (map[uuid.UUID]MediaLocationRef, error)
}

type MediaReader interface {
	MediaAssetExists(ctx context.Context, mediaID uuid.UUID) (bool, error)
}
