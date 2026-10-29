package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/platform/sqlc"
	"messenger/messenger/internal/platform/store"
	"messenger/messenger/internal/shared/app/apperr"
	"strings"
	"uuid"
)

var _ query.Searcher = (*QueryService)(nil)
var _ query.Reader = (*QueryService)(nil)

type QueryService struct {
	*store.Store
}

func NewQueryService(store *store.Store) *QueryService {
	return &QueryService{
		Store: store,
	}
}

func (q *QueryService) Search(ctx context.Context, terms []query.SearchTerm) ([]query.SearchResultItem, error) {
	var qryParts []string
	for _, t := range terms {
		tw := termWildcard(t.Val)
		if len(t.Translits) == 0 {
			qryParts = append(qryParts, tw)
			continue
		}

		tls := make([]string, 0, len(t.Translits))
		for _, tl := range t.Translits {
			tls = append(tls, termWildcard(tl))
		}
		qryParts = append(qryParts, fmt.Sprintf("(%s | %s)", tw, strings.Join(tls, " | ")))
	}

	tsQuery := strings.Join(qryParts, " & ")

	srch, err := q.Queries(ctx).SearchDirectoryEntry(ctx, tsQuery)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result := make([]query.SearchResultItem, 0, len(srch))
	for _, s := range srch {
		result = append(result, query.SearchResultItem{
			ID:          s.TargetID,
			Type:        s.TargetType,
			Handle:      s.Handle,
			DisplayName: s.DisplayName,
		})
	}

	return result, nil
}

func termWildcard(term string) string {
	return fmt.Sprintf("%s:*", term)
}

func (q *QueryService) Get(ctx context.Context, id uuid.UUID, typ entry.TargetType) (query.EntryView, error) {
	entryRow, err := q.Queries(ctx).GetDirectoryEntry(ctx, sqlc.GetDirectoryEntryParams{
		TargetID:   id,
		TargetType: typ.String(),
	})
	if err != nil {
		return query.EntryView{}, apperr.NewNotFound(fmt.Sprintf("Directory entry with id %q and type %q not found", id, typ.String()))
	}

	return query.EntryView{
		ID:          id,
		Type:        typ.String(),
		DisplayName: entryRow.DisplayName,
		AvatarURL:   entryRow.AvatarUrl,
	}, nil
}
