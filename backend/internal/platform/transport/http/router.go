package http

import (
	"fmt"
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/platform/transport/http/contenttype"
	"messenger/messenger/internal/platform/transport/http/headers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	apiHandler StrictServerInterface,
	errorHandler *ErrorHandler,
	tokenService *token.Service,
	wsHandler *WebsocketHandler,
) (http.Handler, error) {
	oas, err := GetSpec()
	if err != nil {
		return nil, fmt.Errorf("failed to load http server specification: %w", err)
	}

	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(AuthMiddleware(tokenService))
	mux.Use(FirewallMiddleware(oas))
	mux.Use(middleware.AllowContentType(contenttype.ApplicationJson))
	mux.Use(middleware.ClientIPFromRemoteAddr)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost:5173", "http://localhost:5173"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{headers.Accept, headers.Authorization, headers.ContentType, headers.XCSRFToken},
		ExposedHeaders:   []string{headers.Link},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/ws", wsHandler.handlerFunc)
	
	strictServer := NewStrictHandlerWithOptions(
		apiHandler, 
		nil, 
		StrictHTTPServerOptions{
			RequestErrorHandlerFunc: errorHandler.Handle,
			ResponseErrorHandlerFunc: errorHandler.Handle,
		},
	)

	return HandlerFromMux(strictServer, mux), nil
}

// import (
// 	"messenger/messenger/internal/auth/infrastructure/token"
// 	"messenger/messenger/internal/platform/transport/http/headers"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/go-chi/cors"
// )

// func routerV1(
// 	tokenService *token.Service,
// 	controller *Controller,
// 	websocketHandler *WebsocketHandler,
// ) *chi.Mux {
// 	r := chi.NewRouter()

// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)
// 	r.Use(middleware.AllowContentType("application/json"))
// 	r.Use(middleware.ClientIPFromRemoteAddr)
// 	r.Use(cors.Handler(cors.Options{
// 		AllowedOrigins:   []string{"https://localhost:5173", "http://localhost:5173"},
// 		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
// 		AllowedHeaders:   []string{headers.Accept, headers.Authorization, headers.ContentType, headers.XCSRFToken},
// 		ExposedHeaders:   []string{headers.Link},
// 		AllowCredentials: true,
// 		MaxAge:           300,
// 	}))
// 	// r.Use(middleware.ClientIPFromXFFTrustedProxies(1))

// 	authMiddleware := AuthMiddleware(tokenService)

// 	registerApi(r, controller, websocketHandler, authMiddleware)

// 	return r
// }

// func registerApi(
// 	r chi.Router,
// 	c *Controller,
// 	ws *WebsocketHandler,
// 	authMiddleware func(http.Handler) http.Handler,
// ) {
// 	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
// 		w.Write([]byte("pong"))
// 	})

// 	// public
// 	r.Group(func(r chi.Router) {
// 		r.Route("/auth", func(r chi.Router) {
// 			r.Post("/register", c.registerUser)
// 			r.Post("/login", c.loginUser)
// 			r.Post("/refresh", c.refreshSession)
// 		})
// 	})

// 	// private
// 	r.Group(func(r chi.Router) {
// 		r.Use(authMiddleware)

// 		r.Get("/ws", ws.handlerFunc)

// 		r.Get("/auth/me", c.getCurrentUser)

// 		r.Route("/me", func(r chi.Router) {
// 			r.Get("/chats", c.getChatList)
// 		})

// 		r.Route("/chats", func(r chi.Router) {
// 			r.Post("/private", c.createPrivateChat)
// 		})

// 		r.Route("/messages", func(r chi.Router) {
// 			r.Post("/", c.sendMessage)
// 		})
// 	})
// }
