package outbox_relay

import (
	"messenger/messenger/internal/bootstrap"
	"messenger/messenger/internal/platform/outbox_relay"

	"github.com/google/wire"
)

type Relay struct {
	OutboxRelay *outbox_relay.OutboxRelay
}

var RelaySet = wire.NewSet(
	wire.Struct(new(Relay), "*"),
	wire.Struct(new(outbox_relay.OutboxRelay), "*"),
	bootstrap.AppSet,
)
