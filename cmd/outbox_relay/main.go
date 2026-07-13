package main

import (
	"context"
	"fmt"
	"messenger/messenger/internal/adapters/out/postgres/queries"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/di"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	ctx := context.Background()

	cfg := config.Init()
	container := di.InitContainer(ctx, cfg)

	eventBus := container.EventBus

	for {
		q := queries.New(container.DB)

		events, err := q.ListUnprocessedEvents(ctx)
		if err != nil {
			return err
		}

		for _, e := range events {
			go eventBus.Dispatch(e)
		}
	}

	return nil
}