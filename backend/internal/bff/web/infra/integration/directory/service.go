package directory

import (
	"context"
	"fmt"
	"messenger/messenger/internal/bff/web/application/port"
	"messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/platform/querybus"
)

var _ port.DirectoryReader = (*Service)(nil)

type Service struct {
	queryBus *querybus.Bus
}

func NewService(
	queryBus *querybus.Bus,
) *Service {
	return &Service{
		queryBus: queryBus,
	}
}

func (s *Service) GetEntries(ctx context.Context, refs []port.EntryRef) ([]port.EntryView, error) {
	qryRefs := make([]query.EntryRef, 0, len(refs))
	for _, r := range refs {
		qryRefs = append(qryRefs, query.EntryRef{
			ID: r.ID,
			Type: r.Type,
		})
	}

	entries, err := s.queryBus.Dispatch[[]query.EntryView](ctx, query.GetEntries{
		Refs: qryRefs,
	})
	if err != nil {
		return nil, fmt.Errorf("get direct entries: %w", err)
	}

	entriesConv := make([]port.EntryView, 0, len(entries))
	for _, e := range entries {
		entriesConv = append(entriesConv, port.EntryView{
			ID: e.ID,
			Type: e.Type,
			DisplayName: e.DisplayName,
			AvatarURL: e.AvatarURL,
		})
	}

	return entriesConv, nil
}