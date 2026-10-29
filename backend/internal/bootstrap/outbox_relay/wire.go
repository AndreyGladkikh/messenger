//go:build wireinject

package outbox_relay

import (
	"messenger/messenger/internal/platform/config"

	"github.com/google/wire"
)

func Initialize(cfg config.Config) (*Relay, func(), error) {
	wire.Build(
		RelaySet,
	)

	return new(Relay), func() {}, nil
}
