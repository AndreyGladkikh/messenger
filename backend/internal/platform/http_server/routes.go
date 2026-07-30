package http_server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func registerApi(mux *chi.Mux, c *Controller) {
	mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.registerUser)
	})

	mux.Route("/chats", func(r chi.Router) {
		r.Post("/private", c.createPrivateChat)
	})

	mux.Route("/messages", func(r chi.Router) {
		r.Post("/", c.sendMessage)
	})
}
