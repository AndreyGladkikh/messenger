package command

import (
	"context"
	"fmt"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/media/domain/model/media"
	"messenger/messenger/internal/media/domain/model/upload"
	"messenger/messenger/internal/shared/app/apperr"
	"messenger/messenger/internal/shared/app/command"
	"uuid"
)

const InitUploadCommandName = "InitUpload"

var _ command.Command = InitUpload{}

type InitUpload struct {
	InitiatorID uuid.UUID
	FileName    string
	Purpose     string
}

func (i InitUpload) IsCommand() {}

func (i InitUpload) Name() string {
	return InitUploadCommandName
}

type InitedUpload struct {
	ID           uuid.UUID
	PresignedURL string
}

type InitUploadHandler struct {
	mediaStorage      port.MediaStorage
	mediaRepo          media.Repository
	uploadRepo         upload.Repository
	mediaLocationStore port.MediaLocationStore
}

func NewInitUploadHandler(
	mediaStorage port.MediaStorage,
	mediaRepo media.Repository,
	uploadRepo upload.Repository,
	mediaLocationStore port.MediaLocationStore,
) *InitUploadHandler {
	return &InitUploadHandler{
		mediaStorage:      mediaStorage,
		mediaRepo:          mediaRepo,
		uploadRepo:         uploadRepo,
		mediaLocationStore: mediaLocationStore,
	}
}

func (h *InitUploadHandler) Handle(ctx context.Context, cmd InitUpload) (InitedUpload, error) {
	m := media.New(cmd.FileName)

	err := h.mediaRepo.Save(ctx, m)
	if err != nil {
		return InitedUpload{}, err
	}

	u := upload.New(
		cmd.InitiatorID,
		m.ID(),
	)

	err = h.uploadRepo.Save(ctx, u)
	if err != nil {
		return InitedUpload{}, err
	}

	var accessLevel port.AccessLevel
	switch cmd.Purpose {
	case "avatar":
		accessLevel = port.AccessLevelPublic
	case "attachment":
		accessLevel = port.AccessLevelPrivate
	default:
		return InitedUpload{}, apperr.NewUnsupportableValue(fmt.Sprintf("Unsupported media purpose %q. Expected one of: %q", cmd.Purpose, []string{"avatar", "attachment"}))
	}

	info, err := h.mediaStorage.PresignUpload(ctx, u.ID(), accessLevel)
	if err != nil {
		return InitedUpload{}, err
	}

	err = h.mediaLocationStore.Save(ctx, m.ID(), info.LocationRef)
	if err != nil {
		return InitedUpload{}, err
	}

	return InitedUpload{
		ID:           u.ID(),
		PresignedURL: info.PresignedURL,
	}, nil
}
