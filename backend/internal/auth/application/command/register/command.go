package register

const CommandName = "RegisterUser"

type Command struct {
	Login string
	Password string
}

func (c *Command) IsCommand() {}

func (c *Command) Name() string {
	return CommandName
}