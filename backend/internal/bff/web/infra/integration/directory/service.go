package directory

import (
	"context"
	"fmt"
	"messenger/messenger/internal/bff/web/app/port"
	drctrQry "messenger/messenger/internal/directory/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
)

var _ port.DirectoryReader = (*Service)(nil)

type Service struct {
	getEntries appQuery.Handler[drctrQry.GetEntries, []drctrQry.EntryView]
}

func NewService(
	getEntries appQuery.Handler[drctrQry.GetEntries, []drctrQry.EntryView],
) *Service {
	return &Service{
		getEntries: getEntries,
	}
}

func (s *Service) GetEntries(ctx context.Context, refs []port.EntryRef) ([]port.EntryView, error) {
	qryRefs := make([]drctrQry.EntryRef, 0, len(refs))
	for _, r := range refs {
		qryRefs = append(qryRefs, drctrQry.EntryRef{
			ID:   r.ID,
			Type: r.Type,
		})
	}

	entries, err := s.getEntries.Handle(ctx, drctrQry.GetEntries{
		Refs: qryRefs,
	})
	if err != nil {
		return nil, fmt.Errorf("get direct entries: %w", err)
	}

	entriesConv := make([]port.EntryView, 0, len(entries))
	for _, e := range entries {
		entriesConv = append(entriesConv, port.EntryView{
			ID:          e.ID,
			Type:        e.Type,
			DisplayName: e.DisplayName,
			AvatarURL:   e.AvatarURL,
		})
	}

	return entriesConv, nil
}
