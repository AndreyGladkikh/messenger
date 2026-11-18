package identity

import (
	"context"
	"messenger/messenger/internal/directory/app/port"
	"messenger/messenger/internal/directory/domain/model"
	idQry "messenger/messenger/internal/identity/app/query"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

var _ port.IdentityReader = (*Service)(nil)

type Service struct {
	getAccount appQuery.Handler[idQry.GetAccount, idQry.AccountView]
}

func NewService(
	getAccount appQuery.Handler[idQry.GetAccount, idQry.AccountView],
) *Service {
	return &Service{
		getAccount: getAccount,
	}
}

func (u *Service) ExistsUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	rspns, err := u.getAccount.Handle(ctx, idQry.GetAccount{ID: userID})
	if err != nil {
		return false, err
	}

	trgtT, err := mapIdentityTypeToTarget(rspns.Type)
	if err != nil {
		return false, err
	}

	if trgtT != model.TargetTypeUser {
		return false, err
	}

	return true, nil
}
