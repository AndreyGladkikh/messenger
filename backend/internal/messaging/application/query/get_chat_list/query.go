package get_chat_list

import "github.com/google/uuid"

const QueryName = "GetChatList"

type Query struct {
	UserID uuid.UUID
}

func (q *Query) IsQuery() {}

func (q *Query) Name() string {
	return QueryName
}
