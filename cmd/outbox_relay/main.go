package main

import (
	"context"
	"log"
	"messenger/messenger/internal/infrastructure/di"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	relay, cleanup, err := di.InitializeOutboxRelay()
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		log.Fatal("event processor: failed to init dependencies", "error", err)
	}

	if err := relay.Run(ctx); err != nil {
		log.Fatal("event processor stopped", "error", err)
	}
}
