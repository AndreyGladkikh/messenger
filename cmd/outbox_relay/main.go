package main

// import (
// 	"context"
// 	"fmt"
// 	"messenger/messenger/internal/infrastructure/db"
// 	"messenger/messenger/internal/infrastructure/config"
// 	"messenger/messenger/internal/infrastructure/di"
// )

// func main() {
// 	if err := run(); err != nil {
// 		fmt.Println(err)
// 	}
// }

// func run() error {
// 	ctx := context.Background()

// 	cfg := config.Load()
// 	container := di.InitContainer(ctx, cfg)

// 	eventBus := container.EventBus

// 	for {
// 		q := queries.New(container.DB)

// 		events, err := q.ListUnprocessedEvents(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		for _, e := range events {
// 			go eventBus.Dispatch(e)
// 		}
// 	}

// 	return nil
// }
//
