package group_conversation

import "uuid"

const CreatedEventName = "conversations.group_conversation.Created.v1"

type Created struct {
	ID       uuid.UUID
	Title    string
	AvatarID *uuid.UUID
}

func (c Created) Name() string {
	return CreatedEventName
}
