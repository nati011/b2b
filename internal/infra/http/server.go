package server

import (
	"context"
	"marketplace/pkg/logger"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

// Options handles server configurations
type Options struct {
	Address      string
	Handler      http.Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// New creates a new HTTP server
func New(opts Options) *Server {
	srv := &http.Server{
		Addr:         opts.Address,
		Handler:      opts.Handler,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
	}

	return &Server{httpServer: srv}
}

// Start begins listening
func (s *Server) Start() error {
	logger.Info("HTTP server starting", "address", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("HTTP server error", "error", err)
		return err
	}
	return nil
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
