package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"messenger/messenger/internal/platform/di"
)

func main() {
	interruptCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	api, cleanup, err := di.InitializeApi()
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		slog.Error("api: failed to init dependencies", "error", err)
		return
	}

	logger := api.Logger
	httpServer := api.HttpServer

	errCh := make(chan error)

	go func() {
		logger.Info("http server started")
		if err := httpServer.Run(); err != nil {
			select {
			case errCh <-err:
				logger.Error("http server failed", "error", err)
			default:
			}
		}
	}()

	select {
	case <-errCh:
	case <-interruptCtx.Done():
	}
	
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpServer.Shutdown(shutdownCtx)
}