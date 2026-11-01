package exchanges

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	NameConversationIntegrationEvents = "conversations.integration-events"
)

func Declare(name string, ch *amqp.Channel) error {
	switch name {
	case NameConversationIntegrationEvents:
		return DeclareConversationIntegrationEvents(ch)
	default:
		return fmt.Errorf("failed to declare an exchange %q: not configured", name)
	}
}

func DeclareConversationIntegrationEvents(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		NameConversationIntegrationEvents,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange %q; %w", NameConversationIntegrationEvents, err)
	}
	return nil
}
