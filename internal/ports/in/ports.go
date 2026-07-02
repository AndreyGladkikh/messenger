package in

import (
	"context"
	"messenger/messenger/internal/application/use_cases/send_private_message"
)

type Ports interface {
	SendPrivateMessage(ctx context.Context, command send_private_message.Command) error
}
