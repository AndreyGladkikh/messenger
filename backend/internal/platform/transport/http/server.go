package http

import (
	"context"
	"errors"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/logger"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

func NewServer(
	cfg *config.Config,
	logger *logger.Logger,
	router http.Handler,
) (*Server, error) {
	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:    cfg.HttpServer.Addr,
			Handler: router,
		},
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Info("http server stopped", "error", ctx.Err())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.httpServer.Shutdown(ctx)
	}()

	s.logger.Info("http server started")
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
