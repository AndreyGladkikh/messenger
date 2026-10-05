package account

import "uuid"

type AccountType string

const (
	AccountTypeUser AccountType = "user"
)

type Account struct {
	ID uuid.UUID
	Type AccountType
}
