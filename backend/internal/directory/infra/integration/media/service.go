package media

import (
	"context"
	"messenger/messenger/internal/directory/app/port"
	"messenger/messenger/internal/media/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

var _ port.MediaReader = (*Service)(nil)

type Service struct {
	mediaExists appQuery.Handler[query.MediaExists, bool]
}

func NewService(
	mediaExists appQuery.Handler[query.MediaExists, bool],
) *Service {
	return &Service{
		mediaExists: mediaExists,
	}
}

func (u *Service) MediaExists(ctx context.Context, mediaID uuid.UUID) (bool, error) {
	exists, err := u.mediaExists.Handle(ctx, query.MediaExists{ID: mediaID})
	if err != nil {
		return false, err
	}

	return exists, nil
}
