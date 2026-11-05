package query

import (
	"context"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

const GetEntriesQueryName = "GetEntries"

var _ query.Query = GetEntries{}

type GetEntries struct {
	Refs []EntryRef
}

func (g GetEntries) Name() string {
	return GetEntriesQueryName
}

type EntryRef struct {
	ID uuid.UUID
	Type string
}

type GetEntriesHandler struct {
	reader Reader
}

func NewGetEntriesHandler(
	reader Reader,
) *GetEntriesHandler {
	return &GetEntriesHandler{
		reader: reader,
	}
}

func (h *GetEntriesHandler) Handle(ctx context.Context, qry GetEntries) ([]EntryView, error) {
	ids := make([]uuid.UUID, 0, len(qry.Refs))
	types := make([]string, 0, len(qry.Refs))
	for _, e := range qry.Refs {
		ids = append(ids, e.ID)
		types = append(types, e.Type)
	}

	entries, err := h.reader.ListByTargets(ctx, ids, types)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
