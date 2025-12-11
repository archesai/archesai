package server

import (
	"github.com/archesai/archesai/pkg/middleware"
)

// Context key re-exports for backwards compatibility.
// These keys are now defined in the middleware package.
const (
	// RequestIDContextKey is the context key for request ID
	RequestIDContextKey = middleware.RequestIDContextKey

	// AuthUserContextKey is the context key for user authentication
	AuthUserContextKey = middleware.AuthUserContextKey

	// AuthClaimsContextKey is the context key for auth claims
	AuthClaimsContextKey = middleware.AuthClaimsContextKey

	// SessionIDContextKey is the context key for session ID
	SessionIDContextKey = middleware.SessionIDContextKey

	// BearerAuthScopes is used by handlers for bearer token authentication
	BearerAuthScopes = middleware.BearerAuthScopes

	// SessionCookieScopes is used by handlers for session cookie authentication
	SessionCookieScopes = middleware.SessionCookieScopes

	// AuthAPIKeyContextKey is the context key for API token
	AuthAPIKeyContextKey = middleware.AuthAPIKeyContextKey

	// AuthMethodContextKey is the context key for the auth method used
	AuthMethodContextKey = middleware.AuthMethodContextKey

	// BearerPrefix is the prefix for Bearer token authentication
	BearerPrefix = middleware.BearerPrefix
)
