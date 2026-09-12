// Package server owns the HTTP server's lifecycle: starting it and shutting
// it down cleanly on SIGINT/SIGTERM, so cmd/api/main.go stays a short boot
// sequence instead of also containing signal-handling and shutdown logic.
package server

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Server wraps a *http.Server built from a Gin engine.
type Server struct {
	httpServer *http.Server
	log        zerolog.Logger
}

// New builds a Server that will listen on ":port" and serve router.
func New(router *gin.Engine, port string, log zerolog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
		log: log,
	}
}

// Run starts the server and blocks until it is shut down, either because
// ListenAndServe returns a non-graceful error, or because the process
// received SIGINT/SIGTERM and the graceful shutdown below completed.
func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		s.log.Info().Str("addr", s.httpServer.Addr).Msg("http server listening")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.log.Info().Msg("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}

	s.log.Info().Msg("server shut down cleanly")
	return nil
}
