package bootstrap

import (
	"messenger/messenger/internal/media/app/command"
	"messenger/messenger/internal/media/app/port"
	"messenger/messenger/internal/media/app/query"
	"messenger/messenger/internal/media/domain/model/media"
	"messenger/messenger/internal/media/domain/model/upload"
	"messenger/messenger/internal/media/infra/object_storage"
	"messenger/messenger/internal/media/infra/persistence"
	"messenger/messenger/internal/media/infra/persistence/repository"
	"messenger/messenger/internal/platform/command_middleware"
	"messenger/messenger/internal/platform/query_middleware"
	appCmd "messenger/messenger/internal/shared/app/command"
	appQry "messenger/messenger/internal/shared/app/query"
	"uuid"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	PersistenceSet,
	CommandHandlersSet,
	QueryHandlersSet,
	object_storage.NewObjectStorage,

	wire.Bind(new(port.MediaStorage), new(*object_storage.ObjectStorage)),
)

var PersistenceSet = wire.NewSet(
	repository.NewMediaRepository,
	repository.NewUploadRepository,
	persistence.NewMediaLocationStore,

	wire.Bind(new(media.Repository), new(*repository.MediaRepository)),
	wire.Bind(new(upload.Repository), new(*repository.UploadRepository)),
	wire.Bind(new(port.MediaLocationStore), new(*persistence.MediaLocationStore)),
)

var CommandHandlersSet = wire.NewSet(
	command.NewComlpeteUploadHandler,
	ProvideComlpeteUploadHandler,
	command.NewInitUploadHandler,
	ProvideInitUploadHandler,
)

func ProvideComlpeteUploadHandler(
	h *command.ComlpeteUploadHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.ComlpeteUpload, uuid.UUID] {
	return command_middleware.Wrap(h, m)
}

func ProvideInitUploadHandler(
	h *command.InitUploadHandler,
	m *command_middleware.Registry,
) appCmd.Handler[command.InitUpload, command.InitedUpload] {
	return command_middleware.Wrap(h, m)
}

var QueryHandlersSet = wire.NewSet(
	query.NewGetDownloadURLHandler,
	ProvideGetMediaDownloadURLHandler,
	query.NewGetPublicURLsHandler,
	ProvideGetPublicURLsHandler,
)

func ProvideGetMediaDownloadURLHandler(
	h *query.GetDownloadURLHandler,
	m *query_middleware.Registry,
) appQry.Handler[query.GetDownloadURL, string] {
	return query_middleware.Wrap(h, m)
}

func ProvideGetPublicURLsHandler(
	h *query.GetPublicURLsHandler,
	m *query_middleware.Registry,
) appQry.Handler[query.GetPublicURLs, map[uuid.UUID]string] {
	return query_middleware.Wrap(h, m)
}
