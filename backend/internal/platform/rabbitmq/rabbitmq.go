package rabbitmq

import (
	"context"
	"encoding/json/jsontext"
	"fmt"
	"messenger/messenger/internal/platform/config"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

func NewConnection(cfg config.Config) (*rmq.AmqpConnection, func(), error) {
	ctx := context.Background()
	env := rmq.NewEnvironment(cfg.Mq.Dsn, nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, func() {
		_ = env.CloseConnections(context.Background())
	}, nil
}

type Message struct {
	Type    string
	Payload jsontext.Value
}

type Consumer interface {
	HandleMessage(ctx context.Context, msg Message) error
}
