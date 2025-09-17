// Package server provides HTTP server infrastructure for the API.
//
// The package includes:
// - HTTP server setup with Echo framework
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

	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/config"
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
	echo   *echo.Echo
	config *config.APIConfig
	logger *slog.Logger
}

// NewServer creates a new API server.
func NewServer(config *config.APIConfig, logger *slog.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	server := &Server{
		echo:   e,
		config: config,
		logger: logger,
	}

	server.SetupMiddleware()
	server.SetupInfrastructureRoutes()

	return server
}

// Echo returns the underlying echo instance for route registration.
func (s *Server) Echo() *echo.Echo {
	return s.echo
}

// ListenAndServe starts the server without signal handling
// This is useful when the caller wants to manage the server lifecycle.
func (s *Server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", int(s.config.Port))
	s.logger.Info("starting server", "address", addr)
	return s.echo.Start(addr)
}

// Start starts the server with built-in signal handling
// This is a convenience method for simple use cases.
func (s *Server) Start() error {
	// Start server in goroutine
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	s.logger.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Error("server forced to shutdown", "error", err)
		return err
	}

	s.logger.Info("server shutdown complete")
	return nil
}

// Shutdown shuts down the server gracefully.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
