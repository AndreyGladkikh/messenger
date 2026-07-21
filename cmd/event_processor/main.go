package main

// import (
// 	"context"
// 	"messenger/messenger/internal/adapters/out/postgres/queries"
// 	"messenger/messenger/internal/platform/config"
// 	"messenger/messenger/internal/platform/di"
// 	"os/signal"
// 	"syscall"
// )

// func main() {
// 	interruptCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
// 	defer stop()

// 	// cfg := config.Load()
// 	// container := di.InitContainer(ctx, cfg)

// 	// eventBus := container.EventBus

// 	// q := queries.New(container.DB)

// 	errCh := make(chan error)

// 	go func() {
// 		for {
// 			events, err := q.ListUnprocessedEvents(ctx)
// 			if err != nil {
// 				errCh <- err
// 			}

// 			for _, e := range events {
// 				go eventBus.Dispatch(e)
// 			}
// 		}
// 	}()

// 	select {
// 	case <-interruptCtx.Done():
// 	}
// }