package server

import (
	"net/http"
)

// Middleware defines the signature for HTTP middleware functions
type Middleware func(http.Handler) http.HandlerFunc

// MiddlewareChain chains multiple middleware functions into a single middleware
func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}
