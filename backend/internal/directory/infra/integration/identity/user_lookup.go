package identity

import (
	"context"
	"messenger/messenger/internal/directory/domain/entry"
	"messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/platform/querybus"
	"uuid"
)

type UserLookup struct {
	qryBus *querybus.Bus
}

func NewUserLookup(
	qryBus *querybus.Bus,
) *UserLookup {
	return &UserLookup{
		qryBus: qryBus,
	}
}

func (u *UserLookup) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	rspns, err := u.qryBus.Dispatch[query.AccountView](ctx, query.GetAccount{ID: id})
	if err != nil {
		return false, err
	}

	trgtT, err := mapIdentityTypeToTarget(rspns.Type)
	if err != nil {
		return false, err
	}

	if trgtT != entry.TargetTypeUser {
		return false, err
	}

	return true, nil
}
