package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/messenger/internal/directory/app/query"
	"messenger/messenger/internal/platform/store"
	"strings"
)

var _ query.Searcher = (*QueryService)(nil)

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
