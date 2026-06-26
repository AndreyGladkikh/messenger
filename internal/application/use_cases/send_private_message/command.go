package send_private_message

import "bytes"

type Command struct {
	RecipientID string
	ChatID string
	MessageBody string
	ResponseToMessageID string
	Attachments []*bytes.Buffer
}

func (c *Command) isCommand() {}

func (c *Command) Name() string {
	return "SendPrivateMessage"
}
