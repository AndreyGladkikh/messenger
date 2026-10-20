//go:build wireinject

package api

import (
	"messenger/messenger/internal/platform/config"

	"github.com/google/wire"
)

func Initialize(cfg config.Config) (*App, func(), error) {
	wire.Build(
		AppSet,
	)

	return new(App), func() {}, nil
}
