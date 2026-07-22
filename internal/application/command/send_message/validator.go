package send_message

import (
	"fmt"
	"messenger/messenger/internal/application/command"
)

type Validator struct {}

// var ErrSenderIDRequired = fmt.Errorf("%w: ", command.ErrValidation)

func ( *Validator) Validate(c *Command) error {
	if c.SenderID == "" {
		return fmt.Errorf("%w: SenderID required", command.ErrValidation)
	}
	if c.ChatID == "" {
		return fmt.Errorf("%w: ChatID required", command.ErrValidation)
	}
	if c.MessageBody == "" {
		return fmt.Errorf("%w: MessageBody required", command.ErrValidation)
	}

	return nil
}