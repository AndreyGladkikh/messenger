package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	initial()
}

func initial() {
	
}

func async() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	
	defer func() {
		fmt.Println("deferred exec")
	}()

	go func () {
		fmt.Println("go started")

		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			log.Fatal(err)
		}

		var shouldSendSigEveryTick bool
		sendSig := func() {
			fmt.Println("sig emitted")
			if err := p.Signal(os.Interrupt); err != nil {
				log.Fatal(err)
			}
		}

		ticks := time.Tick(1*time.Second)
		timeout := time.After(3*time.Second)

		select {
		case <-ticks:
			fmt.Println("tick")
			if shouldSendSigEveryTick {
				sendSig()
			}
		case <-timeout:
			fmt.Println("timeout")
			sendSig()
			shouldSendSigEveryTick = true
		}
	}()

	
	fmt.Println("blocked")
	<-ctx.Done()
	fmt.Println("resumed")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fmt.Println("ctx with timeout created")

	go func () {
		fmt.Println("go 2 started")

		ticks := time.Tick(1*time.Second)

		select {
		case <-ticks:
			fmt.Println("tick 2")
		}
	}()

	
	fmt.Println("timeout blocked")
	<-ctx.Done()
	fmt.Println("timeout resumed")

	c := 3
	_ = c
	fmt.Println("end")
}