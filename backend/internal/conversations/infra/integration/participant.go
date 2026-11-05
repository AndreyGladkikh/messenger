package integration

import (
	"context"
	"errors"
	idQry "messenger/messenger/internal/identity/app/query"
	"messenger/messenger/internal/shared/app/apperr"
	appQuery "messenger/messenger/internal/shared/app/query"
	"uuid"
)

type ParticipantLookup struct {
	getAccount appQuery.Handler[idQry.GetAccount, idQry.AccountView]
}

func NewParticipantLookup(
	getAccount appQuery.Handler[idQry.GetAccount, idQry.AccountView],
) *ParticipantLookup {
	return &ParticipantLookup{
		getAccount: getAccount,
	}
}

func (p *ParticipantLookup) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := p.getAccount.Handle(ctx, idQry.GetAccount{ID: id})
	if err != nil {
		if err, ok := errors.AsType[*apperr.Error](err); ok && err.Kind == apperr.KindNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
