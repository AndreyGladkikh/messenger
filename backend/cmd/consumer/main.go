package main

import (
	"context"
	"flag"
	"log"
	"messenger/messenger/internal/bootstrap/consumer"
	"messenger/messenger/internal/platform/config"
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

	if runner, ok := container.Consumers[cnsmrName]; ok {
		if err := runner.Run(ctx); err != nil {
			return err
		}
	}

	// runner := getRunner(container.Consumers, cnsmrName)
	// if runner != nil {
	// 	if err := runner.Run(ctx); err != nil {
	// 		return err
	// 	}
	// }

	return nil

	// // todo bootstrap
	// consCfg, ok := NewConsumerConfigs(Consumers{})[cnsmrName]
	// if !ok {
	// 	return fmt.Errorf("not found config for consumer %q", cnsmrName)
	// }

	// rmqConn := container.RMQConn

	// err = prepareConsumer(ctx, consCfg, rmqConn)
	// if err != nil {
	// 	return err
	// }

	// consumer, err := rmqConn.NewConsumer(ctx, consCfg.qName, nil)
	// if err != nil {
	// 	log.Panicf("Failed to create consumer: %v", err)
	// }
	// defer func() { _ = consumer.Close(context.Background()) }()

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// for {
	// 	delivery, err := consumer.Receive(ctx)
	// 	if err != nil {
	// 		if errors.Is(err, context.Canceled) {
	// 			return nil
	// 		}
	// 		return fmt.Errorf("Failed to receive a message %w", err)
	// 	}
	// 	msg := delivery.Message()
	// 	var body string
	// 	if len(msg.Data) > 0 {
	// 		body = string(msg.Data[0])
	// 	}
	// 	log.Printf(" [x] %s", body)
	// 	err = delivery.Accept(ctx)
	// 	if err != nil {
	// 		return fmt.Errorf("Failed to accept message %w", err)
	// 	}
	// }
}
