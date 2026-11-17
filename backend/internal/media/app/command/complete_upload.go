package command

import (
	"context"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/media/domain/model/media"
	"messenger/messenger/internal/media/domain/model/upload"
	"messenger/messenger/internal/shared/app/command"
	"uuid"
)

const ComlpeteUploadCommandName = "ComlpeteUpload"

var _ command.Command = ComlpeteUpload{}

type ComlpeteUpload struct {
	UploadID uuid.UUID
}

func (c ComlpeteUpload) IsCommand() {}

func (c ComlpeteUpload) Name() string {
	return ComlpeteUploadCommandName
}

type ComlpeteUploadHandler struct {
	mediaStorage      port.MediaStorage
	mediaLocationStore port.MediaLocationStore
	mediaRepo          media.Repository
	uploadRepo         upload.Repository
}

func NewComlpeteUploadHandler(
	mediaStorage port.MediaStorage,
	mediaLocationStore port.MediaLocationStore,
	mediaRepo media.Repository,
	uploadRepo upload.Repository,
) *ComlpeteUploadHandler {
	return &ComlpeteUploadHandler{
		mediaStorage:      mediaStorage,
		mediaLocationStore: mediaLocationStore,
		mediaRepo:          mediaRepo,
		uploadRepo:         uploadRepo,
	}
}

func (h *ComlpeteUploadHandler) Handle(ctx context.Context, cmd ComlpeteUpload) (uuid.UUID, error) {
	u, err := h.uploadRepo.Get(ctx, cmd.UploadID)
	if err != nil {
		return uuid.UUID{}, err
	}

	objRef, err := h.mediaLocationStore.Get(ctx, u.MediaID)
	if err != nil {
		return uuid.UUID{}, err
	}

	_, err = h.mediaStorage.GetObjectInfo(ctx, objRef)
	if err != nil {
		return uuid.UUID{}, err
	}

	u.Complete()

	err = h.uploadRepo.Save(ctx, u)
	if err != nil {
		return uuid.UUID{}, err
	}

	m, err := h.mediaRepo.Get(ctx, u.MediaID)
	if err != nil {
		return uuid.UUID{}, err
	}

	m.MarkReady()

	err = h.mediaRepo.Save(ctx, m)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.MediaID, nil
}
