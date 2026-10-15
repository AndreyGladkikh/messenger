package rabbitmq

import (
	"context"
	"fmt"
	"messenger/messenger/internal/platform/config"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

func NewConnection(cfg *config.Config) (*rmq.AmqpConnection, func(), error) {
	ctx := context.Background()
	env := rmq.NewEnvironment(cfg.RabbitMq.Addr, nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, func() {
		_ = env.CloseConnections(context.Background())
	}, nil
}