package e2e

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"mime"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMediaUploadDownload_Success(t *testing.T) {
	server, _ := initApi(t)

	alice := alice(t, server)

	storageClient := &http.Client{Timeout: 10 * time.Second}

	photo := image.NewRGBA(image.Rect(0, 0, 2, 2))
	var pngBytes, jpegBytes bytes.Buffer
	require.NoError(t, png.Encode(&pngBytes, photo))
	require.NoError(t, jpeg.Encode(&jpegBytes, photo, nil))

	for _, tc := range []struct {
		name              string
		content           []byte
		uploadContentType string
		contentType       string
		purpose       string
	}{
		{"photo.png", pngBytes.Bytes(), "application/octet-stream", "image/png", "avatar"},
		{"Фото с отпуска 🏖.JPG", jpegBytes.Bytes(), "", "image/jpeg", "attachment"},
		{`notes "final"; 100% + #1.txt`, []byte("original file contents\n"), "", "text/plain; charset=utf-8", "avatar"},
		{"photo", pngBytes.Bytes(), "image/png", "image/png", "attachment"},
		{"data.unknown-media-format", []byte{0, 1, 2, 255}, "application/x-custom", "application/x-custom", "avatar"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response, _ := alice.httpClient.initMediaUpload(t, appHttp.InitMediaUploadJSONRequestBody{FileName: tc.name, Purpose: tc.purpose})
			require.Equal(t, http.StatusCreated, response.StatusCode)
			upload := decode[appHttp.InitedMediaUpload](t, response.Body)

			request, err := http.NewRequestWithContext(t.Context(), http.MethodPut, upload.PresignedUrl, bytes.NewReader(tc.content))
			require.NoError(t, err)
			if tc.uploadContentType != "" {
				request.Header.Set("Content-Type", tc.uploadContentType)
			}
			response, err = storageClient.Do(request)
			require.NoError(t, err)
			response.Body.Close()
			require.Equal(t, http.StatusOK, response.StatusCode)

			response, _ = alice.httpClient.completeMediaUpload(t, appHttp.CompleteMediaUploadRequestObject{UploadId: upload.Id})
			require.Equal(t, http.StatusOK, response.StatusCode)
			asset := decode[appHttp.Media](t, response.Body)

			response, _ = alice.httpClient.getMediaDownloadURL(t, appHttp.GetDownloadMediaUrlRequestObject{MediaId: asset.MediaId})
			require.Equal(t, http.StatusOK, response.StatusCode)
			download := decode[appHttp.DownloadMediaUrl](t, response.Body)

			request, err = http.NewRequestWithContext(t.Context(), http.MethodGet, download.Url, nil)
			require.NoError(t, err)
			response, err = storageClient.Do(request)
			require.NoError(t, err)
			defer response.Body.Close()
			require.Equal(t, http.StatusOK, response.StatusCode)

			disposition, parameters, err := mime.ParseMediaType(response.Header.Get("Content-Disposition"))
			require.NoError(t, err)
			assert.Equal(t, "attachment", disposition)
			assert.Equal(t, tc.name, parameters["filename"])
			assert.Equal(t, tc.contentType, response.Header.Get("Content-Type"))
			content, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Equal(t, tc.content, content)
		})
	}
}
