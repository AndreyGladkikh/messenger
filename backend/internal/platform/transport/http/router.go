package http

import (
	"fmt"
	"messenger/messenger/internal/identity/infra/token"
	"messenger/messenger/internal/platform/transport/http/contenttype"
	"messenger/messenger/internal/platform/transport/http/headers"
	"messenger/messenger/internal/platform/transport/websocket"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	apiHandler StrictServerInterface,
	errorHandler *ErrorHandler,
	tokenService *token.Service,
	wsHandler *websocket.Handler,
) (*chi.Mux, error) {
	oas, err := GetSpec()
	if err != nil {
		return nil, fmt.Errorf("failed to load http server specification: %w", err)
	}

	mux := chi.NewMux()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost:5173", "http://localhost:5173", "localhost:5173"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{headers.Accept, headers.Authorization, headers.ContentType, headers.XCSRFToken},
		ExposedHeaders:   []string{headers.Link},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Logger)
	mux.Use(AuthMiddleware(tokenService))
	mux.Use(middleware.AllowContentType(contenttype.ApplicationJson))
	mux.Use(middleware.ClientIPFromRemoteAddr)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)

	mux.Group(func(r chi.Router) {
		r.Use(RequestValidator(oas, errorHandler))

		HandlerFromMux(NewStrictHandlerWithOptions(
			apiHandler,
			nil,
			StrictHTTPServerOptions{
				RequestErrorHandlerFunc:  errorHandler.Handle,
				ResponseErrorHandlerFunc: errorHandler.Handle,
			},
		), r)
	})

	mux.Get("/ws", wsHandler.ServeHTTP)

	// ###
	// todo rm
	ah := apiHandler.(*ApiHandler)
	mux.Post("/test", ah.Test)
	// ###

	return mux, nil
}
