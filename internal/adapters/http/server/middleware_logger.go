package server

import (
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// loggerMiddleware logs HTTP requests
func (s *Server) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		requestID, _ := r.Context().Value(RequestIDContextKey).(string)

		s.logger.Info("request",
			"id", requestID,
			"method", r.Method,
			"uri", r.RequestURI,
			"status", wrapped.statusCode,
			"latency", duration,
			"remote_ip", r.RemoteAddr,
		)
	})
}
