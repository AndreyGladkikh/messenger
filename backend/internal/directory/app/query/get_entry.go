package query

import (
	"context"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/shared/app/query"
	"uuid"
)

const GetEntryQueryName = "GetEntry"

var _ query.Query = GetEntry{}

type GetEntry struct {
	ID   uuid.UUID
	Type string
}

func (g GetEntry) Name() string {
	return GetEntryQueryName
}

var _ query.Handler[GetEntry, EntryView] = (*GetEntryHandler)(nil)

type GetEntryHandler struct {
	Reader Reader
}

func NewGetEntryHandler(reader Reader) *GetEntryHandler {
	return &GetEntryHandler{
		Reader: reader,
	}
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
