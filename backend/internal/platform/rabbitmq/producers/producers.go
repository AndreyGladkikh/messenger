package producers

import (
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/platform/rabbitmq/exchanges"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producers struct {
	ConversationIntegrationEventsProducer *ConversationIntegrationEventsProducer
}

type ConversationIntegrationEventsProducer struct {
	*rabbitmq.Producer
}

func NewConversationIntegrationEventsProducer(
	logger *logger.Logger,
	amqpConn *amqp.Connection,
) (*ConversationIntegrationEventsProducer, func(), error) {
	p, cleanup, err := rabbitmq.NewProducer(logger, amqpConn, exchanges.NameConversationIntegrationEvents)
	if err != nil {
		return nil, nil, err
	}
	return &ConversationIntegrationEventsProducer{p}, cleanup, nil
}
