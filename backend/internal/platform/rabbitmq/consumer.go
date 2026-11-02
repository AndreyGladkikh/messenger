package rabbitmq

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"log"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/rabbitmq/exchanges"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerRunner interface {
	Run(ctx context.Context) error
}

var _ ConsumerRunner = (*Consumer)(nil)

type Consumer struct {
	logger  *logger.Logger
	cfg     ConsumerConfig
	name    string
	rmqChan *amqp.Channel
	queue   amqp.Queue
	handler Handler
}

func NewConsumer(
	logger *logger.Logger,
	rmqConn *amqp.Connection,
	handler Handler,
	qName string,
	cfg ConsumerConfig,
) (*Consumer, func(), error) {
	ch, err := rmqConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("rabbitmq consumer: failed to open a channel: %w", err)
	}
	cleanup := func() {
		ch.Close()
	}

	q, err := ch.QueueDeclare(
		qName,
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		return nil, cleanup, fmt.Errorf("rabbitmq consumer: failed to declare a queue %q: %w", qName, err)
	}

	for _, x := range cfg.Exchanges {
		err := exchanges.Declare(x.Name, ch)
		if err != nil {
			return nil, cleanup, err
		}

		if len(x.BindingKeys) == 0 {
			x.BindingKeys = append(x.BindingKeys, "*")
		}

		for _, bk := range x.BindingKeys {
			err = ch.QueueBind(
				q.Name,
				bk,
				x.Name,
				false,
				nil,
			)
			if err != nil {
				return nil, cleanup, fmt.Errorf("failed to bind a queue %q: %w", qName, err)
			}
		}
	}

	return &Consumer{
		logger:  logger,
		rmqChan: ch,
		queue:   q,
		handler: handler,
	}, cleanup, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	msgs, err := c.rmqChan.Consume(
		c.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq consumer: failed to register a consumer %q: %w", c.name, err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			logErr := func(msg string, args ...any) {
				args = append([]any{
					"consumer", c.name,
					"msg", d.Body,
				}, args...)
				c.logger.ErrorContext(ctx, fmt.Sprintf("rabbitmq consumer: %s", msg), args...)
			}

			log.Printf("Received a message: %s", d.Body)

			var envlp Envelope
			err := json.Unmarshal(d.Body, &envlp)
			if err != nil {
				logErr("failed to decode message", "err", err)
				d.Nack(false, false)
				continue
			}

			err = c.handler.Handle(ctx, envlp)
			if err != nil {
				logErr("failed to handle message", "err", err)
				d.Nack(false, false)
				continue
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

type Handler interface {
	Handle(ctx context.Context, msg Envelope) error
}

type ConsumerConfig struct {
	Exchanges []XConfig
	QName     string
}

type XConfig struct {
	Name        string
	BindingKeys []string
}

// func NewConsumerConfigs(consumers Consumers) ConsumerConfigs {
// 	return ConsumerConfigs{
// 		rabbitmq.QNameRealtimeBroadcastDelivery: {
// 			exchanges: []xConfig{
// 				{
// 					name: rabbitmq.XNameConversationEvents,
// 				},
// 			},
// 			qName:       rabbitmq.QNameRealtimeBroadcastDelivery,
// 			handlerFunc: consumers.RealTimeDeliveryPeers.Handle,
// 		},
// 		rabbitmq.QNameRealtimePersonalDelivery: {
// 			exchanges: []xConfig{
// 				{
// 					name:        rabbitmq.XNameConversationEvents,
// 					bindingKeys: []string{"direct-started.v1"},
// 				},
// 			},
// 			qName: rabbitmq.QNameRealtimePersonalDelivery,
// 		},
// 	}
// }
