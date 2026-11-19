package e2e

import (
	"bytes"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"
)

type uploadedMedia struct {
	id        uuid.UUID
	publicURL *string
}

func uploadMedia(t *testing.T, u *user, fileName string, file []byte, purpose string) uploadedMedia {
	t.Helper()

	response, _ := u.httpClient.initMediaUpload(t, appHttp.InitMediaUploadJSONRequestBody{
		FileName: fileName,
		Purpose:  purpose,
	})
	require.Equal(t, http.StatusCreated, response.StatusCode)
	upload := decode[appHttp.InitedMediaUpload](t, response.Body)

	request, err := http.NewRequestWithContext(t.Context(), http.MethodPut, upload.PresignedUrl, bytes.NewReader(file))
	require.NoError(t, err)
	response, err = (&http.Client{Timeout: 10 * time.Second}).Do(request)
	require.NoError(t, err)
	response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)

	response, _ = u.httpClient.completeMediaUpload(t, appHttp.CompleteMediaUploadRequestObject{UploadId: upload.Id})
	require.Equal(t, http.StatusOK, response.StatusCode)
	asset := decode[appHttp.Media](t, response.Body)

	var publicURL *string
	switch purpose {
	case "avatar":
		url := testCDNHost + "/public/" + upload.Id.String() + "/original"
		publicURL = &url
	}
	
	return uploadedMedia{
		id: asset.MediaId,
		publicURL: publicURL,
	}
}
