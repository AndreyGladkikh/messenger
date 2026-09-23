package http

import (
	"context"
	"errors"
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/logger"
	// "messenger/messenger/internal/platform/transport/http/api"
	"net/http"
	"time"

	// "github.com/go-chi/chi/v5"
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

	// apiServer *api.Server,
) *Server {
	router := router(
		tokenService,
		controller,
		websocketHandler,
	)

	// h := api.HandlerFromMux(api.NewStrictHandler(apiServer, nil), chi.NewMux())

	server := &http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
		// Handler: h,
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
