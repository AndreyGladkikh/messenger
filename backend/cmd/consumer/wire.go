//go:build wireinject

package main

import (
	"github.com/google/wire"
)

func initContainer() (*container, func(), error) {
	wire.Build(containerSet)

	return new(container), func() {}, nil
}
