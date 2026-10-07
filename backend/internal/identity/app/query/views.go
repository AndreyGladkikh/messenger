package query

import "uuid"

type AccountView struct {
	ID   uuid.UUID
	Type string
}
