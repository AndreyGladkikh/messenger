// +build wireinject

package di

import (
	"messenger/messenger/internal/infrastructure/event"

	"github.com/google/wire"
)

func InitializeApi() (*Api, func(), error) {
	wire.Build(
		ApiSet,
	)

	return new(Api), func() {}, nil
}

func InitializeEventProcessor() (*event.Processor, func(), error) {
	wire.Build(
		OutboxRelaySet,
	)

	return new(event.Processor), func() {}, nil
}
