package object_storage

import (
	"context"
	"messenger/messenger/internal/media/app/port"
	platformS3 "messenger/messenger/internal/platform/s3"
	"mime"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPresignDownloadPreservesFileNameAndType(t *testing.T) {
	client := awsS3.NewFromConfig(aws.Config{
		Region: "us-east-1",
		Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "test", SecretAccessKey: "test-secret"}, nil
		}),
	}, func(options *awsS3.Options) {
		options.BaseEndpoint = aws.String("https://storage.example.com")
		options.UsePathStyle = true
	})
	storage := NewObjectStorage(platformS3.NewPresigner(client), client)
	ref := port.MediaLocationRef{Bucket: "media", Key: "opaque-upload-id"}

	for _, tc := range []struct {
		name        string
		contentType string
	}{
		{"photo.png", "image/png"},
		{"holiday photo.JPG", "image/jpeg"},
		{"animation.gif", "image/gif"},
		{"report.pdf", "application/pdf"},
		{"notes.txt", "text/plain; charset=utf-8"},
		{"Фото с отпуска 🏖.png", "image/png"},
		{`photo "final"; 100% + #1.png`, "image/png"},
		{"photo\r\nX-Injected: value.png", "image/png"},
		{"archive.unknown-media-format", ""},
		{"no-extension", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			info, err := storage.PresignDownload(t.Context(), ref, tc.name)
			require.NoError(t, err)
			parsed, err := url.Parse(info.PresignedURL)
			require.NoError(t, err)
			assert.Equal(t, "/media/opaque-upload-id", parsed.Path)

			params := parsed.Query()
			header := params.Get("response-content-disposition")
			disposition, attributes, err := mime.ParseMediaType(header)
			require.NoError(t, err)
			assert.Equal(t, "attachment", disposition)
			assert.Equal(t, tc.name, attributes["filename"])
			assert.NotContains(t, header, "\r")
			assert.NotContains(t, header, "\n")
			assert.Equal(t, tc.contentType, params.Get("response-content-type"))
			assert.Equal(t, tc.contentType != "", params.Has("response-content-type"))
			assert.NotEmpty(t, params.Get("X-Amz-Signature"))
			assert.Equal(t, "60", params.Get("X-Amz-Expires"))
		})
	}
}
