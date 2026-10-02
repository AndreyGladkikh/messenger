package list_user_chats

import "uuid"

const QueryName = "GetChatList"

type Query struct {
	UserID uuid.UUID
}

func (q *Query) IsQuery() {}

func (q *Query) Name() string {
	return QueryName
}
