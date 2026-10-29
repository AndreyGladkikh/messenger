//go:build wireinject

package api

import (
	"messenger/messenger/internal/platform/config"

	"github.com/google/wire"
)

func Initialize(cfg config.Config) (*Api, func(), error) {
	wire.Build(
		ApiSet,
	)

	return new(Api), func() {}, nil
}
