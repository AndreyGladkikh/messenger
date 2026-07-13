package main

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5"

	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/di"
	"messenger/messenger/internal/platform/http_server"
	// "net/http"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	// commandBus := new(command.Bus)
	// commandBus.Register()

	ctx := context.Background()

	cfg := config.Init()
	container := di.InitContainer(ctx, cfg)

	commandBus := commandbus.BuildCommandBus(
		container.TxManager,
		container.SendMessageHandler,
	)

	httpServer := http_server.NewServer(commandBus)

	if err := httpServer.Run(); err != nil {
		return err
	}

	return nil

	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func (w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })
	// return http.ListenAndServe(":3000", r)
}
