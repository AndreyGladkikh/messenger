package http_server

import (
	"context"
	"errors"
	"messenger/messenger/internal/platform/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server     *http.Server
}

func NewServer(
	cfg *config.Config,
	controller *Controller,
) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	registerApi(r, controller)

	server := &http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: r,
	}

	return &Server{
		server:     server,
	}
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}