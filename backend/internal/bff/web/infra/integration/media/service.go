package media

import (
	"context"
	"fmt"
	"messenger/messenger/internal/bff/web/app/port"
	"messenger/messenger/internal/media/app/query"
	appQry "messenger/messenger/internal/shared/app/query"
	"uuid"
)

var _ port.MediaReader = (*Service)(nil)

type Service struct {
	getMediaPublicURLs appQry.Handler[query.GetPublicURLs, map[uuid.UUID]string]
}

func NewService(
	getMediaPublicURLs appQry.Handler[query.GetPublicURLs, map[uuid.UUID]string],
) *Service {
	return &Service{
		getMediaPublicURLs: getMediaPublicURLs,
	}
}

func (s *Service) GetAvatarURLs(ctx context.Context, avatarIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	urls, err := s.getMediaPublicURLs.Handle(ctx, query.GetPublicURLs{MediaIDs: avatarIDs})
	if err != nil {
		return nil, fmt.Errorf("get media public urls")
	}

	return urls, nil
}
