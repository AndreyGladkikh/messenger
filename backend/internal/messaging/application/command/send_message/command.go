package send_message

import (
	"bytes"

	"github.com/google/uuid"
)

const CommandName = "SendMessage"

type Command struct {
	SenderID         uuid.UUID
	ChatID           uuid.UUID
	MessageBody      string
	ReplyToMessageID uuid.UUID
	Attachments      []*bytes.Buffer
}

func (c *Command) IsCommand() {}

func (c *Command) Name() string {
	return CommandName
}
