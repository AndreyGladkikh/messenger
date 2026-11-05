package identity

import (
	"context"
	"messenger/messenger/internal/directory/domain/entry"
	idQry "messenger/messenger/internal/identity/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

type UserLookup struct {
	getAccount appQuery.Handler[idQry.GetAccount, idQry.AccountView]
}

func NewUserLookup(
	getAccount appQuery.Handler[idQry.GetAccount, idQry.AccountView],
) *UserLookup {
	return &UserLookup{
		getAccount: getAccount,
	}
}

func (u *UserLookup) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	rspns, err := u.getAccount.Handle(ctx, idQry.GetAccount{ID: id})
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
