package main

import (
	"context"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"messenger/messenger/internal/bootstrap/api"
	"messenger/messenger/internal/platform/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	api, cleanup, err := api.Initialize(cfg)
	if err != nil {
		slog.Error("api: failed to init dependencies", "error", err)
		return
	}
	defer cleanup()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Go(func() {
		errCh <- api.HttpServer.Run(ctx)
	})

	interruptCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg.Go(func() {
		select {
		case <-interruptCtx.Done():
		case <-errCh:
		}
		cancel()
	})

	wg.Wait()
}
