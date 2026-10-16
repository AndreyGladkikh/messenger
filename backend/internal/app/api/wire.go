//go:build wireinject

package api

import (
	"github.com/google/wire"
)

func Initialize() (*App, func(), error) {
	wire.Build(
		AppSet,
	)

	return new(App), func() {}, nil
}
