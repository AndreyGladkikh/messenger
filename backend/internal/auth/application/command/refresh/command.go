package refresh

const CommandName = "Refresh"

type Command struct {
	RefreshToken string
}

func (c *Command) IsCommand() {}

func (c *Command) Name() string {
	return CommandName
}
