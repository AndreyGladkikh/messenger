package consumers

import (
	cnvEvent "messenger/messenger/internal/conversations/app/event"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/platform/rabbitmq/exchanges"
	"messenger/messenger/internal/platform/realtime"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ConsumerNameRealtimePersonalDelivery = "realtime:personal-delivery"
)

type Consumers struct {
	RealtimePersonalDeliveryConsumer *RealtimePersonalDeliveryConsumer
}

type RealtimePersonalDeliveryConsumer struct {
	*rabbitmq.Consumer
}

func NewRealtimePersonalDeliveryConsumer(
	logger *logger.Logger,
	rmqConn *amqp.Connection,
	handler *realtime.PersonalDeliveryHandler,
) (*RealtimePersonalDeliveryConsumer, func(), error) {
	cfg := rabbitmq.ConsumerConfig{
		Exchanges: []rabbitmq.XConfig{
			{
				Name:        exchanges.NameConversationIntegrationEvents,
				BindingKeys: []string{cnvEvent.DirectConversationStartedEventName},
			},
		},
		QName: rabbitmq.QNameRealtimePersonalDelivery,
	}

	c, cleanup, err := rabbitmq.NewConsumer(logger, rmqConn, handler, rabbitmq.QNameRealtimePersonalDelivery, cfg)
	if err != nil {
		return nil, nil, err
	}
	return &RealtimePersonalDeliveryConsumer{c}, cleanup, nil
}
