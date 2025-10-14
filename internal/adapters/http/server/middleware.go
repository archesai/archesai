package server

import (
	"net/http"
)

// WrapHandler applies all middleware to a handler
func (s *Server) WrapHandler(h http.Handler) http.Handler {
	// Apply middleware in reverse order (last middleware wraps first)
	h = s.timeoutMiddleware(h)
	h = s.rateLimitMiddleware(h)
	h = s.securityMiddleware(h)
	h = s.corsMiddleware(h)
	h = s.recoverMiddleware(h)
	h = s.loggerMiddleware(h)
	h = s.requestIDMiddleware(h)
	return h
}
