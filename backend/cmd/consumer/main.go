package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"messenger/messenger/internal/bootstrap/consumer"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/realtime/messaging"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	ctx := context.Background()

	cfg := config.Load()

	cnsmrName := flag.Args()[0]

	container, cleanup, err := consumer.Initialize(cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	// todo bootstrap
	consCfg, ok := NewConsumerConfigs(Consumers{})[cnsmrName]
	if !ok {
		return fmt.Errorf("not found config for consumer %q", cnsmrName)
	}

	rmqConn := container.RMQConn

	err = prepareConsumer(ctx, consCfg, rmqConn)
	if err != nil {
		return err
	}

	consumer, err := rmqConn.NewConsumer(ctx, consCfg.qName, nil)
	if err != nil {
		log.Panicf("Failed to create consumer: %v", err)
	}
	defer func() { _ = consumer.Close(context.Background()) }()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	for {
		delivery, err := consumer.Receive(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return fmt.Errorf("Failed to receive a message %w", err)
		}
		msg := delivery.Message()
		var body string
		if len(msg.Data) > 0 {
			body = string(msg.Data[0])
		}
		log.Printf(" [x] %s", body)
		err = delivery.Accept(ctx)
		if err != nil {
			return fmt.Errorf("Failed to accept message %w", err)
		}
	}
}

const (
	consumerNameRealtimeDeliveryPeers = "realtime:delivery:peers"
	consumerNameRealtimeDeliveryUsers = "realtime:delivery:users"
)

type consumerConfig struct {
	exchanges   []xConfig
	qName       string
	handlerFunc func(ctx context.Context, msg []byte) error
}

type xConfig struct {
	name        string
	bindingKeys []string
}

type HandlerFunc func(ctx context.Context, msg []byte) error

type ConsumerConfigs map[string]consumerConfig

// var consumerConfigs = map[string]consumerConfig {
// 	consumerNameRealtimeDeliveryPeers: {
// 		exchanges: []xConfig{
// 			{
// 				name: "conversation.events",
// 				bindingKeys: []string{"peer:*"},
// 			},
// 		},
// 		qName: "realtime.deliver_peers",
// 	},
// 	consumerNameRealtimeDeliveryUsers: {
// 		exchanges: []xConfig{
// 			{
// 				name: "conversation.events",
// 				bindingKeys: []string{"user:*"},
// 			},
// 		},
// 		qName: "realtime.deliver_users",
// 	},
// }

type Consumers struct {
	RealTimeDeliveryPeers *messaging.Peers
}

func NewConsumerConfigs(consumers Consumers) ConsumerConfigs {
	return ConsumerConfigs{
		consumerNameRealtimeDeliveryPeers: {
			exchanges: []xConfig{
				{
					name:        "conversation.events",
					bindingKeys: []string{"peer.*"},
				},
			},
			qName:       "realtime.deliver_peers",
			handlerFunc: consumers.RealTimeDeliveryPeers.Handle,
		},
		consumerNameRealtimeDeliveryUsers: {
			exchanges: []xConfig{
				{
					name:        "conversation.events",
					bindingKeys: []string{"user.*"},
				},
			},
			qName: "realtime.deliver_users",
		},
	}
}

func prepareConsumer(ctx context.Context, cfg consumerConfig, conn *rmq.AmqpConnection) error {
	for _, x := range cfg.exchanges {
		_, err := conn.Management().DeclareExchange(ctx, &rmq.TopicExchangeSpecification{Name: x.name})
		if err != nil {
			return fmt.Errorf("failed to declare an exchange: %w", err)
		}
	}

	qInfo, err := conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{
		Name: cfg.qName,
	})
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	for _, x := range cfg.exchanges {
		for _, bk := range x.bindingKeys {
			_, err = conn.Management().Bind(ctx, &rmq.ExchangeToQueueBindingSpecification{
				SourceExchange:   x.name,
				DestinationQueue: qInfo.Name(),
				BindingKey:       bk,
			})
			if err != nil {
				return fmt.Errorf("failed to bind a queue: %w", err)
			}
		}
	}

	return nil
}
