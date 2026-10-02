package create_private_chat

import "uuid"

var CommandName = "CreatePrivateChat"

type Command struct {
	InitiatorID    uuid.UUID
	ChatWithUserID uuid.UUID
}

func (c *Command) IsCommand() {}

func (c *Command) Name() string {
	return "CreatePrivateChat"
}
