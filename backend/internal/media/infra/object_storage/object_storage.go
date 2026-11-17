package object_storage

import (
	"context"
	"fmt"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/platform/s3"
	"mime"
	"path/filepath"
	"uuid"

	awsClient "github.com/aws/aws-sdk-go-v2/service/s3"
)

var _ port.MediaStorage = (*ObjectStorage)(nil)

type ObjectStorage struct {
	presigner *s3.Presigner
	client    *awsClient.Client
}

func NewObjectStorage(
	presigner *s3.Presigner,
	client *awsClient.Client,
) *ObjectStorage {
	return &ObjectStorage{
		presigner: presigner,
		client:    client,
	}
}

func (o *ObjectStorage) PresignUpload(ctx context.Context, uploadID uuid.UUID, accessLevel port.AccessLevel) (port.UploadInfo, error) {
	bucket := s3.BucketPrivate
	switch accessLevel {
	case port.AccessLevelPrivate:
		bucket = s3.BucketPrivate
	case port.AccessLevelPublic:
		bucket = s3.BucketPublic
	}

	key := fmt.Sprintf("%s/original", uploadID.String())

	req, err := o.presigner.PutObject(ctx, bucket, key, 60)
	if err != nil {
		return port.UploadInfo{}, err
	}

	return port.UploadInfo{
		PresignedURL: req.URL,
		LocationRef: port.MediaLocationRef{
			Bucket: bucket,
			Key:    key,
		},
	}, nil
}

func (o *ObjectStorage) GetObjectInfo(ctx context.Context, ref port.MediaLocationRef) (port.MediaInfo, error) {
	_, err := o.client.HeadObject(ctx, &awsClient.HeadObjectInput{
		Bucket: &ref.Bucket,
		Key:    &ref.Key,
	})
	if err != nil {
		return port.MediaInfo{}, fmt.Errorf("head object: %w", err)
	}

	return port.MediaInfo{}, nil
}

func (o *ObjectStorage) PresignDownload(ctx context.Context, ref port.MediaLocationRef, fileName string) (port.DownloadInfo, error) {
	resp, err := o.presigner.GetObject(ctx, ref.Bucket, ref.Key, 60, s3.GetObjectOptions{
		ContentDisposition: mime.FormatMediaType("attachment", map[string]string{"filename": fileName}),
		ContentType:        mime.TypeByExtension(filepath.Ext(fileName)),
	})
	if err != nil {
		return port.DownloadInfo{}, err
	}

	return port.DownloadInfo{
		PresignedURL: resp.URL,
	}, nil
}
