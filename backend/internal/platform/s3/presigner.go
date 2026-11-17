package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct {
	PresignClient *s3.PresignClient
}

type GetObjectOptions struct {
	ContentDisposition string
	ContentType        string
}

func NewPresigner(s3Client *s3.Client) *Presigner {
	return &Presigner{s3.NewPresignClient(s3Client)}
}

func (presigner *Presigner) GetObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64, options GetObjectOptions) (*v4.PresignedHTTPRequest, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	if options.ContentDisposition != "" {
		input.ResponseContentDisposition = aws.String(options.ContentDisposition)
	}
	if options.ContentType != "" {
		input.ResponseContentType = aws.String(options.ContentType)
	}

	request, err := presigner.PresignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		return nil, fmt.Errorf("get a presigned request to get %v:%v: %w", bucketName, objectKey, err)
	}

	return request, nil
}

func (presigner *Presigner) PutObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		return nil, fmt.Errorf("get a presigned request to put %v:%v: %w", bucketName, objectKey, err)
	}

	return request, nil
}
