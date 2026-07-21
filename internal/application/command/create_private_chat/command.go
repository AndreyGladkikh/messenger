package create_private_chat

type Command struct {
	InitiatorID string
	ChatWithUserID string
}

func (c *Command) IsCommand() {}

func (c *Command) Name() string {
	return "CreatePrivateChat"
}
