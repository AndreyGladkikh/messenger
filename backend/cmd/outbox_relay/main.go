package main

import (
	"context"
	"log"
	"messenger/messenger/internal/bootstrap/outbox_relay"
	"messenger/messenger/internal/platform/config"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	relay, cleanup, err := outbox_relay.Initialize(cfg)
	if err != nil {
		log.Fatal("event processor: failed to init dependencies", "error", err)
	}
	defer cleanup()

	if err := relay.OutboxRelay.Run(ctx); err != nil {
		log.Fatal("event processor stopped", "error", err)
	}
}
