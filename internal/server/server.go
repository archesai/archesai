// Package server provides HTTP server infrastructure for the API.
//
// The package includes:
// - HTTP server setup with Echo framework
// - Middleware configuration (CORS, logging, recovery)
// - WebSocket support for real-time communication
// - OpenAPI documentation serving
// - Graceful shutdown handling
package server

import "time"

// Server configuration constants
const (
	// DefaultPort is the default server port
	DefaultPort = "8080"

	// DefaultReadTimeout is the default read timeout
	DefaultReadTimeout = 30 * time.Second

	// DefaultWriteTimeout is the default write timeout
	DefaultWriteTimeout = 30 * time.Second

	// DefaultIdleTimeout is the default idle timeout
	DefaultIdleTimeout = 120 * time.Second

	// DefaultShutdownTimeout is the timeout for graceful shutdown
	DefaultShutdownTimeout = 10 * time.Second

	// DefaultMaxHeaderBytes is the maximum header size
	DefaultMaxHeaderBytes = 1 << 20 // 1 MB
)

// WebSocket constants
const (
	// WebSocketReadBufferSize is the WebSocket read buffer size
	WebSocketReadBufferSize = 1024

	// WebSocketWriteBufferSize is the WebSocket write buffer size
	WebSocketWriteBufferSize = 1024

	// WebSocketHandshakeTimeout is the WebSocket handshake timeout
	WebSocketHandshakeTimeout = 10 * time.Second

	// WebSocketPingPeriod is the WebSocket ping period
	WebSocketPingPeriod = 54 * time.Second

	// WebSocketPongTimeout is the WebSocket pong timeout
	WebSocketPongTimeout = 60 * time.Second
)

// Middleware priority constants (lower number = higher priority)
const (
	// MiddlewarePriorityRecover runs first to catch panics
	MiddlewarePriorityRecover = 1

	// MiddlewarePriorityLogger logs all requests
	MiddlewarePriorityLogger = 2

	// MiddlewarePriorityCORS handles CORS headers
	MiddlewarePriorityCORS = 3

	// MiddlewarePriorityAuth handles authentication
	MiddlewarePriorityAuth = 10
)
