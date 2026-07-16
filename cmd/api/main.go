package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"messenger/messenger/internal/platform/di"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	interruptCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	api, cleanup, err := di.InitializeApi()
	if err != nil {
		return err
	}
	defer cleanup()

	logger := api.Logger
	httpServer := api.HttpServer

	errCh := make(chan error)

	go func() {
		if err := httpServer.Run(); err != nil {
			errCh <-err
		}
	}()

	select {
	case err := <-errCh:
		fmt.Printf("errCh: %s", err.Error())
		logger.Error("api error: %v", err)
	case <-interruptCtx.Done():
		fmt.Println("interruptCtx")
	}
	

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpServer.Shutdown(shutdownCtx)
	fmt.Printf("shutdown")

	return nil
}
