package get_current_user

import "uuid"

const QueryName = "GetCurrentUser"

type Query struct {
	UserID uuid.UUID
}

func (q Query) IsQuery() {}

func (q Query) Name() string {
	return QueryName
}
