package lookup

import (
	"context"
	"errors"
	"messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/platform/querybus"
	"messenger/messenger/internal/shared/domain"
	"uuid"
)

type ParticipantLookup struct {
	qryBus *querybus.Bus
}

func NewParticipantLookup(
	qryBus *querybus.Bus,
) *ParticipantLookup {
	return &ParticipantLookup{
		qryBus: qryBus,
	}
}

func (p *ParticipantLookup) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := p.qryBus.Dispatch[query.AccountView](ctx, query.GetAccount{ID: id})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
