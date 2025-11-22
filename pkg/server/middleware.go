package server

import (
	"net/http"
)

// WrapHandler applies all middleware to a handler
func (s *APIServer) WrapHandler(h http.Handler) http.Handler {
	// Apply middleware in reverse order (last middleware wraps first)
	h = TimeoutMiddleware(h)
	h = RateLimitMiddleware(h)
	h = SecurityMiddleware(h)
	h = CorsMiddleware(h, s.config.Cors)
	h = RecoverMiddleware(h)
	h = LoggerMiddleware(h)
	h = RequestIDMiddleware(h)
	return h
}
