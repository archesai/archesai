package middleware

import (
	"net/http"
)

// Middleware defines the signature for HTTP middleware functions
type Middleware func(http.Handler) http.HandlerFunc

// Chain chains multiple middleware functions into a single middleware
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}
