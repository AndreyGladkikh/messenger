package main

import (
	"context"
	"log"
	"messenger/messenger/internal/bootstrap/outbox_relay"
	"messenger/messenger/internal/platform/config"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	relay, cleanup, err := outbox_relay.Initialize(cfg)
	if err != nil {
		log.Fatal("outbox relay: failed to init dependencies", "err", err)
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errc := make(chan error, 1)
	go func() {
		if err := relay.OutboxRelay.Run(ctx); err != nil {
			errc <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("outbox relay: shutting down")
	case err := <-errc:
		log.Fatal("outbox relay: failed to run", "err", err)
	}

	cancel()

	<-time.After(30 * time.Second)
}
