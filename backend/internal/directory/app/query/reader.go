package query

import (
	"context"
	"messenger/messenger/internal/directory/domain/entry"
	"uuid"
)

type Reader interface {
	Get(ctx context.Context, id uuid.UUID, typ entry.TargetType) (EntryView, error)
}
