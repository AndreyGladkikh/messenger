package create_chat

type Command struct {
	Type     string
	ChatName string
}

func (c *Command) Name() string {
	return "CreateChat"
}
