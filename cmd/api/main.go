package main

import (
	"fmt"
	"messenger/messenger/internal/application/command"
	"messenger/messenger/internal/platform/http_server"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	commandBus := new(command.Bus)
	commandBus.Register()

	httpServer := http_server.NewServer()

	if err := httpServer.Run(); err != nil {
		return err
	}

	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func (w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })
	// return http.ListenAndServe(":3000", r)
}
