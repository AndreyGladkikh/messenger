package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	ctx := context.Background()

	container, cleanup, err := initContainer()
	if err != nil {
		return err
	}
	defer cleanup()

	cnsmrName := consumerNameDirectoryUpdateSearchIndex

	consCfg, ok := consumerConfigs[cnsmrName]
	if !ok {
		return fmt.Errorf("not found config for consumer %q", cnsmrName)
	}

	rmqConn := container.rmqConn

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
	consumerNameDirectoryUpdateSearchIndex = "directory:update_search"
)

type consumerConfig struct {
	xNames []string
	qName string
	handlerFunc func(ctx context.Context, msg []byte) error
}

var consumerConfigs = map[string]consumerConfig {
	consumerNameDirectoryUpdateSearchIndex: {
		xNames: []string{
			"directory.entry_registered",
			"directory.entry_updated",
		},
		qName: "directory.update_search",
	},
}

func prepareConsumer(ctx context.Context, cfg consumerConfig, conn *rmq.AmqpConnection) error {
	for _, xName := range cfg.xNames {
		_, err := conn.Management().DeclareExchange(ctx, &rmq.FanOutExchangeSpecification{Name: xName})
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

	for _, xName := range cfg.xNames {
		_, err = conn.Management().Bind(ctx, &rmq.ExchangeToQueueBindingSpecification{
			SourceExchange:   xName,
			DestinationQueue: qInfo.Name(),
			BindingKey:       "",
		})
		if err != nil {
			return fmt.Errorf("failed to bind a queue: %w", err)
		}
	}

	return nil
}
