package producers

import (
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/platform/rabbitmq/exchanges"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producers struct {
	ConversationEventsProducer *ConversationEventsProducer
}

type ConversationEventsProducer struct {
	*rabbitmq.Producer
}

func NewConversationEventsProducer(
	logger *logger.Logger,
	amqpConn *amqp.Connection,
) (*ConversationEventsProducer, func(), error) {
	p, cleanup, err := rabbitmq.NewProducer(logger, amqpConn, exchanges.NameConversationIntegrationEvents)
	if err != nil {
		return nil, nil, err
	}
	return &ConversationEventsProducer{p}, cleanup, nil
}
