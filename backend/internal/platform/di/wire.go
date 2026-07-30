// +build wireinject

package di

import (
	"messenger/messenger/internal/messaging/infrastructure/outbox_relay"

	"github.com/google/wire"
)

func InitializeApi() (*Api, func(), error) {
	wire.Build(
		ApiSet,
	)

	return new(Api), func() {}, nil
}

func InitializeOutboxRelay() (*outbox_relay.OutboxRelay, func(), error) {
	wire.Build(
		OutboxRelaySet,
	)

	return new(outbox_relay.OutboxRelay), func() {}, nil
}
