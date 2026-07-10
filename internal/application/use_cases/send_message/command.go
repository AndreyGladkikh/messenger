package send_message

import "bytes"

type Command struct {
	SenderID string
	ChatID string
	MessageBody string
	ReplyToMessageID string
	Attachments []*bytes.Buffer
}

func (c *Command) IsCommand() {}

func (c *Command) Name() string {
	return "SendMessage"
}
