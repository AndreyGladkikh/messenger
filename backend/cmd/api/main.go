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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.Load()

	api, cleanup, err := api.Initialize(cfg)
	if err != nil {
		slog.Error("api: failed to init dependencies", "error", err)
		return
	}
	defer cleanup()

	var wg sync.WaitGroup
	errCh := make(chan error)

	wg.Go(func() {
		if err := api.HttpServer.Run(ctx); err != nil {
			errCh <- err
		}
	})

	wg.Go(func() {
		<-errCh
		cancel()
	})

	wg.Wait()
}
