package http_server

import (
	"context"
	"errors"
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/logger"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	logger *logger.Logger
}

func NewServer(
	cfg *config.Config,
	logger *logger.Logger,
	controller *Controller,
	websocketHandler *WebsocketHandler,
	tokenService *token.Service,
) *Server {
	router := router(
		tokenService,
		controller,
		websocketHandler,
	)

	server := &http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
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
