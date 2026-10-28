package query

import (
	"context"
	"uuid"
)

type Reader interface {
	GetAccount(context.Context, uuid.UUID) (AccountView, error)
	GetUserIDByLogin(context.Context, string) (uuid.UUID, error)
}
