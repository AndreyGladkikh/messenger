package main

import (
	"context"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"messenger/messenger/internal/infrastructure/di"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	api, cleanup, err := di.InitializeApi()
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		slog.Error("api: failed to init dependencies", "error", err)
		return
	}

	httpServer := api.HttpServer

	var wg sync.WaitGroup
	errCh := make(chan error)

	wg.Go(func() {
		if err := httpServer.Run(ctx); err != nil {
			errCh <- err
		}
	})

	select {
	case <-errCh:
		cancel()
	default:
		wg.Wait()
	}
}
