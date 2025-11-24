package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/archesai/archesai/pkg/auth"
)

// Auth context keys
const (
	// AuthUserContextKey is the context key for user authentication
	AuthUserContextKey contextKey = "auth.user"

	// AuthClaimsContextKey is the context key for auth claims
	AuthClaimsContextKey contextKey = "auth.claims"

	// SessionIDContextKey is the context key for session ID
	SessionIDContextKey contextKey = "auth.sessionID"

	// BearerAuthScopes is used by handlers for bearer token authentication
	BearerAuthScopes = "bearerAuth.Scopes"

	// SessionCookieScopes is used by handlers for session cookie authentication
	SessionCookieScopes = "sessionCookie.Scopes"

	// AuthAPIKeyContextKey is the context key for API token
	AuthAPIKeyContextKey = "auth_api_token"

	// AuthMethodContextKey is the context key for the auth method used
	AuthMethodContextKey = "auth_method"

	// BearerPrefix is the prefix for Bearer token authentication
	BearerPrefix = "Bearer"
)

// Claims represents JWT claims with user information
type Claims struct {
	UserID         uuid.UUID `json:"user_id"`
	SessionID      uuid.UUID `json:"session_id"`
	Email          string    `json:"email"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	Roles          []string  `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware provides authentication middleware using the auth service.
type AuthMiddleware struct {
	authService *auth.Service
}

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(authService *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth creates middleware that validates JWT tokens.
func (am *AuthMiddleware) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response := NewUnauthorizedResponse("missing authorization header", r.URL.Path)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Check for Bearer token format
			tokenString := strings.TrimPrefix(authHeader, BearerPrefix+" ")
			if tokenString == authHeader {
				response := NewUnauthorizedResponse(
					"invalid authorization header format",
					r.URL.Path,
				)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Validate access token using auth service
			claims, err := am.authService.ValidateAccessToken(tokenString)
			if err != nil {
				response := NewUnauthorizedResponse("invalid token", r.URL.Path)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), AuthUserContextKey, claims.UserID)
			ctx = context.WithValue(ctx, AuthClaimsContextKey, claims)
			ctx = context.WithValue(ctx, SessionIDContextKey, claims.SessionID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth creates middleware for optional authentication.
func (am *AuthMiddleware) OptionalAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No auth header, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			// Check for Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != BearerPrefix {
				// Invalid format, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			tokenString := parts[1]

			// Try to validate token, but don't fail if invalid
			claims, err := am.authService.ValidateAccessToken(tokenString)
			if err == nil && claims != nil {
				// Add claims to context
				ctx := context.WithValue(r.Context(), AuthUserContextKey, claims.UserID)
				ctx = context.WithValue(ctx, AuthClaimsContextKey, claims)
				ctx = context.WithValue(ctx, SessionIDContextKey, claims.SessionID)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClaimsFromContext retrieves the claims from the context
func GetClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(AuthClaimsContextKey).(*Claims)
	return claims, ok
}

// RequireRole creates middleware that requires specific roles
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaimsFromContext(r.Context())
			if !ok {
				response := NewUnauthorizedResponse("no authentication claims", r.URL.Path)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Check if user has any of the required roles
			for _, requiredRole := range roles {
				for _, userRole := range claims.Roles {
					if userRole == requiredRole {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			response := NewForbiddenResponse("insufficient permissions", r.URL.Path)
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		})
	}
}
