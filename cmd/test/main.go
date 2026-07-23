package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	channels()
}

func channels() {
	// ch := make(chan bool)

	var wg sync.WaitGroup

	wg.Go(func() {
		fmt.Println("go 1: started")
		ticker := time.Tick(1*time.Second)
		after := time.After(2*time.Second)

		for {
			select {
			case <-ticker:
				fmt.Println("go 1: tick")
			case <-after:
				fmt.Println("go 1: end")
				return
			default:
				fmt.Println("go 1: default")
				time.Sleep(500*time.Millisecond)
			}
		}
	})

	wg.Go(func() {
		fmt.Println("go 2: started")
		ticker := time.Tick(2*time.Second)
		after := time.After(6*time.Second)

		for {
			select {
			case <-ticker:
				fmt.Println("go 2: tick")
			case <-after:
				fmt.Println("go 2: end")
				return
			default:
				fmt.Println("go 2: default")
				time.Sleep(500*time.Millisecond)
			}
		}
	})

	fmt.Println("main: start waiting")
	wg.Wait()
	fmt.Println("main: end waiting")
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