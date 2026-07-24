package main

import (
	"context"
	"log"
	"messenger/messenger/internal/infrastructure/di"
	"os/signal"
	"syscall"
)

func main() {
	interruptCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	processor, cleanup, err := di.InitializeEventProcessor()
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		log.Fatal("event processor: failed to init dependencies", "error", err)
		return
	}

	processor.Run(interruptCtx)

	// errCh := make(chan error)

	// go func() {
	// 	processor.Run(interruptCtx)
	// }()

	// select {
	// case <-interruptCtx.Done():
	// }
}
