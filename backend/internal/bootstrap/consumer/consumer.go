package consumer

import (
	"messenger/messenger/internal/bootstrap"

	"github.com/google/wire"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type Consumer struct {
	RMQConn *rmq.AmqpConnection
}

var ApiSet = wire.NewSet(
	wire.Struct(new(Consumer), "*"),
	bootstrap.AppSet,
)
