package query

import (
	"context"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

var _ query.Query = GetEntry{}

type GetEntry struct {
	ID   uuid.UUID
	Type string
}

func (g GetEntry) IsQuery() {}

func (g GetEntry) Name() string {
	return "GetEntry"
}

var _ query.Handler[GetEntry, EntryView] = (*GetEntryHandler)(nil)

type GetEntryHandler struct {
	Reader Reader
}

func (h *GetEntryHandler) Handle(ctx context.Context, qry GetEntry) (EntryView, error) {
	typ, err := entry.ParseTargetType(qry.Type)
	if err != nil {
		return EntryView{}, err
	}

	view, err := h.Reader.Get(ctx, qry.ID, typ)
	if err != nil {
		return EntryView{}, err
	}

	return view, nil
}
