package login

import "net/netip"

const CommandName = "LoginUser"

type Command struct {
	Login     string
	Password  string
	UserAgent string
	IP        netip.Addr
}

func (c *Command) IsCommand() {}

func (c *Command) Name() string {
	return CommandName
}
