package http_server

import (
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/platform/http_server/headers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func router(
	tokenService *token.Service,
	controller *Controller,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost:5173", "http://localhost:5173"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{headers.Accept, headers.Authorization, headers.ContentType, headers.XCSRFToken},
		ExposedHeaders:   []string{headers.Link},
		AllowCredentials: true,
		MaxAge: 300,
	}))
	// r.Use(middleware.ClientIPFromXFFTrustedProxies(1))

	registerApi(r, tokenService, controller)

	return r
}

func registerApi(r chi.Router, tokenService *token.Service, c *Controller) {
	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("pong"))
	})

	// public
	r.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", c.registerUser)
			r.Post("/login", c.loginUser)
			r.Post("/refresh", c.refreshSession)
		})
	})

	// private
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(tokenService))

		r.Route("/me", func(r chi.Router) {
			r.Get("/chats", c.getChatList)
		})

		r.Route("/chats", func(r chi.Router) {
			r.Post("/private", c.createPrivateChat)
		})

		r.Route("/messages", func(r chi.Router) {
			r.Post("/", c.sendMessage)
		})
	})
}
