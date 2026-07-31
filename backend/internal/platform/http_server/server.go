package http_server

import (
	"context"
	"errors"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server *http.Server
	logger *logger.Logger
}

func NewServer(
	cfg *config.Config,
	logger *logger.Logger,
	controller *Controller,
) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(AuthMiddleware)
	r.Use(middleware.ClientIPFromRemoteAddr)
	// r.Use(middleware.ClientIPFromXFFTrustedProxies(1))

	registerApi(r, controller)

	server := &http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: r,
	}

	return &Server{
		server: server,
		logger: logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Info("http server stopped", "error", ctx.Err())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.server.Shutdown(ctx)
	}()

	s.logger.Info("http server started")
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
