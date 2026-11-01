package consumer

import (
	"messenger/messenger/internal/bootstrap"
	"messenger/messenger/internal/platform/rabbitmq/consumers"
	"messenger/messenger/internal/platform/realtime"

	"github.com/google/wire"
)

type Consumer struct {
	Consumers *consumers.Consumers
}

var ConsumerSet = wire.NewSet(
	wire.Struct(new(Consumer), "*"),
	realtime.NewPersonalDeliveryHandler,
	bootstrap.AppSet,
)
