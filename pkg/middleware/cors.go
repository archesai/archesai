package middleware

import (
	"net/http"
	"strings"
)

// CORS creates CORS middleware with the given allowed origins
func CORS(origins string) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		return cors(next, origins)
	}
}

// cors handles CORS
func cors(next http.Handler, origins string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originList := strings.Split(origins, ",")
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, o := range originList {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().
			Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
