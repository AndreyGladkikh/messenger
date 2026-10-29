//go:build wireinject

package consumer

import (
	"messenger/messenger/internal/platform/config"

	"github.com/google/wire"
)

func Initialize(cfg config.Config) (*Consumer, func(), error) {
	wire.Build(
		ApiSet,
	)

	return new(Consumer), func() {}, nil
}
