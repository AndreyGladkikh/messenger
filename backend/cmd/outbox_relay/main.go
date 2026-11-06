package main

import (
	"context"
	"log"
	"messenger/messenger/internal/platform/bootstrap/outbox_relay"
	"messenger/messenger/internal/platform/config"
	"os/signal"
	"syscall"
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
		errc <- relay.OutboxRelay.Run(ctx)
	}()

	interruptCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-interruptCtx.Done()
		cancel()
	}()

	<-errc
}
