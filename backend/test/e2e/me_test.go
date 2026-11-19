package e2e

import (
	"bytes"
	"image"
	"image/png"
	appHttp "messenger/messenger/internal/platform/transport/http"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMyProfile_Success(t *testing.T) {
	server, _ := initApi(t)
	user := newUser(t, server, "alice")

	storageClient := &http.Client{Timeout: 10 * time.Second}

	handle := "alice_handle"
	displayName := "Alice Display Name"

	response, _ := user.httpClient.initMediaUpload(t, appHttp.InitMediaUploadJSONRequestBody{FileName: "alice_avatar", Purpose: "avatar"})
	require.Equal(t, http.StatusCreated, response.StatusCode)
	upload := decode[appHttp.InitedMediaUpload](t, response.Body)

	photo := image.NewRGBA(image.Rect(0, 0, 2, 2))
	var pngBytes bytes.Buffer
	require.NoError(t, png.Encode(&pngBytes, photo))

	request, err := http.NewRequestWithContext(t.Context(), http.MethodPut, upload.PresignedUrl, bytes.NewReader(pngBytes.Bytes()))
	require.NoError(t, err)
	response, err = storageClient.Do(request)
	require.NoError(t, err)
	response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)

	response, _ = user.httpClient.completeMediaUpload(t, appHttp.CompleteMediaUploadRequestObject{UploadId: upload.Id})
	require.Equal(t, http.StatusOK, response.StatusCode)
	mediaAsset := decode[appHttp.Media](t, response.Body)

	resp, err := user.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle:      &handle,
		DisplayName: displayName,
		AvatarMediaId: &mediaAsset.MediaId,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, _ = user.httpClient.getMyProfile(t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	profile := decode[appHttp.Profile](t, resp.Body)
	assert.Equal(t, &handle, profile.Handle)
	assert.Equal(t, displayName, profile.DisplayName)
}

func TestCreateMyProfile_HandleAlreadyExists(t *testing.T) {
	server, _ := initApi(t)
	aliceA := newUser(t, server, "aliceA")
	aliceB := newUser(t, server, "aliceB")

	aliceAHandle := "alice"
	aliceBHandle := "alice"

	resp, _ := aliceA.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle:      &aliceAHandle,
		DisplayName: "Alice A",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, _ = aliceB.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle:      &aliceBHandle,
		DisplayName: "Alice B",
	})
	require.Equal(t, http.StatusConflict, resp.StatusCode)
	problem := decode[appHttp.ProblemDetails](t, resp.Body)
	assert.Equal(t, problemTypesRegistry.AlreadyExists, *problem.Type)

	resp, _ = aliceB.httpClient.getMyProfile(t)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	aliceBHandle = "alice_b"
	resp, _ = aliceB.httpClient.createMyProfile(t, appHttp.CreateMyProfileJSONRequestBody{
		Handle:      &aliceBHandle,
		DisplayName: "Alice B",
	})
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateMyProfile_ProfileAlreadyExists(t *testing.T) {
	server, _ := initApi(t)
	user := newUser(t, server, "alice")

	handle := "alice"
	displayName := "Alice"

	request := appHttp.CreateMyProfileJSONRequestBody{
		Handle:      &handle,
		DisplayName: displayName,
	}

	resp, _ := user.httpClient.createMyProfile(t, request)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, _ = user.httpClient.createMyProfile(t, request)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, _ = user.httpClient.getMyProfile(t)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	stored := decode[appHttp.Profile](t, resp.Body)
	assert.Equal(t, &handle, stored.Handle)
	assert.Equal(t, displayName, stored.DisplayName)
}
