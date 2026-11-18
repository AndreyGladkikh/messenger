package directory

import (
	"context"
	"fmt"
	"messenger/messenger/internal/bff/web/app/port"
	directoryQry "messenger/messenger/internal/directory/app/query"
	directoryView "messenger/messenger/internal/directory/app/view"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

var _ port.DirectoryReader = (*Service)(nil)

type Service struct {
	getUserProfiles appQuery.Handler[directoryQry.GetUserProfiles, []directoryView.UserProfile]
	searchDirectoryEntry appQuery.Handler[directoryQry.Search, []directoryQry.SearchResultItem]
}

func NewService(
	getUserProfiles appQuery.Handler[directoryQry.GetUserProfiles, []directoryView.UserProfile],
	searchDirectoryEntry appQuery.Handler[directoryQry.Search, []directoryQry.SearchResultItem],
) *Service {
	return &Service{
		getUserProfiles: getUserProfiles,
		searchDirectoryEntry: searchDirectoryEntry,
	}
}

func (s *Service) GetUserProfiles(ctx context.Context, ids []uuid.UUID) ([]port.UserProfileView, error) {
	ups, err := s.getUserProfiles.Handle(ctx, directoryQry.GetUserProfiles{UserIDs: ids})
	if err != nil {
		return nil, fmt.Errorf("get user profiles: %w", err)
	}

	profView := make([]port.UserProfileView, 0, len(ups))
	for _, up := range ups {
		profView = append(profView, port.UserProfileView{
			Handle: up.Handle,
			Bio: up.Bio,
			Birthday: up.Birthday,
			DisplayName: up.DisplayName,
			AvatarID:    up.AvatarID,
			UserID: up.UserID,
		})
	}

	return profView, nil
}

func (s *Service) Search(ctx context.Context, str string) ([]port.SearchResultItem, error) {
	searchResults, err := s.searchDirectoryEntry.Handle(ctx, directoryQry.Search{Str: str})
	if err != nil {
		return nil, fmt.Errorf("search directory entry: %w", err)
	}

	items := make([]port.SearchResultItem, 0, len(searchResults))
	for _, r := range searchResults {
		items = append(items, port.SearchResultItem{
			ID: r.ID,
			Type: r.Type,
			Handle: r.Handle,
			DisplayName: r.DisplayName,
			AvatarID: r.AvatarID,
		})
	}

	return items, nil
}
