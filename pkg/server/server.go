// Package server provides HTTP server infrastructure for the API.
//
// The package includes:
// - HTTP server setup with Go 1.22+ stdlib routing
// - Middleware configuration (CORS, logging, recovery)
// - WebSocket support for real-time communication
// - OpenAPI documentation serving
// - Graceful shutdown handling
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/archesai/archesai/pkg/config"
)

// Server configuration constants.
const (
	// DefaultPort is the default server port.
	DefaultPort = "8080"

	// DefaultReadTimeout is the default read timeout.
	DefaultReadTimeout = 30 * time.Second

	// DefaultWriteTimeout is the default write timeout.
	DefaultWriteTimeout = 30 * time.Second

	// DefaultIdleTimeout is the default idle timeout.
	DefaultIdleTimeout = 120 * time.Second

	// DefaultShutdownTimeout is the timeout for graceful shutdown.
	DefaultShutdownTimeout = 10 * time.Second

	// DefaultMaxHeaderBytes is the maximum header size.
	DefaultMaxHeaderBytes = 1 << 20 // 1 MB
)

// WebSocket constants.
const (
	// WebSocketReadBufferSize is the WebSocket read buffer size.
	WebSocketReadBufferSize = 1024

	// WebSocketWriteBufferSize is the WebSocket write buffer size.
	WebSocketWriteBufferSize = 1024

	// WebSocketHandshakeTimeout is the WebSocket handshake timeout.
	WebSocketHandshakeTimeout = 10 * time.Second

	// WebSocketPingPeriod is the WebSocket ping period.
	WebSocketPingPeriod = 54 * time.Second

	// WebSocketPongTimeout is the WebSocket pong timeout.
	WebSocketPongTimeout = 60 * time.Second
)

// Middleware priority constants (lower number = higher priority).
const (
	// MiddlewarePriorityRecover runs first to catch panics.
	MiddlewarePriorityRecover = 1

	// MiddlewarePriorityLogger logs all requests.
	MiddlewarePriorityLogger = 2

	// MiddlewarePriorityCORS handles CORS headers.
	MiddlewarePriorityCORS = 3

	// MiddlewarePriorityAuth handles authentication.
	MiddlewarePriorityAuth = 10
)

// Server represents the API server.
type Server struct {
	mux    *http.ServeMux
	server *http.Server
	config *config.API
}

// NewServer creates a new API server.
func NewServer(config *config.API) *Server {
	mux := http.NewServeMux()

	server := &Server{
		mux:    mux,
		config: config,
	}

	return server
}

// Mux returns the underlying http.ServeMux for route registration.
func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

// ListenAndServe starts the server without signal handling
// This is useful when the caller wants to manage the server lifecycle.
func (s *Server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", int(s.config.Port))

	s.server = &http.Server{
		Addr:           addr,
		Handler:        s.mux,
		ReadTimeout:    DefaultReadTimeout,
		WriteTimeout:   DefaultWriteTimeout,
		IdleTimeout:    DefaultIdleTimeout,
		MaxHeaderBytes: DefaultMaxHeaderBytes,
	}

	slog.Info("starting server", "address", addr)
	return s.server.ListenAndServe()
}

// Start starts the server with built-in signal handling
// This is a convenience method for simple use cases.
func (s *Server) Start() error {
	// Start server in goroutine
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		return err
	}

	slog.Info("server shutdown complete")
	return nil
}

// Shutdown shuts down the server gracefully.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}
