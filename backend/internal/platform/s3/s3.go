package s3

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/config"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	BucketMedia = "media"
	BucketPublic = "public"
	BucketPrivate = "private"

	// PrefixMedia = "media"
	// PrefixAvatars = "avatars"
)

func New(cfg config.Config) (*s3.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

    awsCfg, err := awsCfg.LoadDefaultConfig(ctx)
    if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
    }

    client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.ObjectStorage.Dsn)
		o.UsePathStyle = true
	})

	return client, nil
}