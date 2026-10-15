package main

import (
	"messenger/messenger/internal/platform/di"
	"messenger/messenger/internal/platform/rabbitmq"

	"github.com/google/wire"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type container struct {
	rmqConn *rmq.AmqpConnection
}

func newContainer(
	rmqConn *rmq.AmqpConnection,
) *container {
	return &container{
		rmqConn: rmqConn,
	}
}

var containerSet = wire.NewSet(
	newContainer,
	di.CommonSet,
	rabbitmq.NewConnection,
)
