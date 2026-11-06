package consumer

import (
	"messenger/messenger/internal/bootstrap"
	"messenger/messenger/internal/platform/rabbitmq"
	"messenger/messenger/internal/platform/realtime"

	"github.com/google/wire"
)

type Consumer struct {
	Consumers map[string]rabbitmq.ConsumerRunner
}

var ConsumerSet = wire.NewSet(
	wire.Struct(new(Consumer), "*"),
	bootstrap.AppSet,
	realtime.NewPersonalDeliveryHandler,
)
