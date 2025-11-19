package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// recoverMiddleware recovers from panics
func (s *Server) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := r.Context().Value(RequestIDContextKey).(string)
				slog.Error("panic recovered",
					"id", requestID,
					"error", err,
				)

				response := NewInternalServerErrorResponse(
					"An unexpected error occurred",
					r.URL.Path,
				)

				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
