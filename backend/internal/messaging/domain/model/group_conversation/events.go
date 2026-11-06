package group_conversation

import "uuid"

const CreatedEventName = "conversations.group_conversation.Created.v1"

type Created struct {
	ID uuid.UUID
	Title string
	AvatarURL *string
}

func (c Created) Name() string {
	return CreatedEventName
}
