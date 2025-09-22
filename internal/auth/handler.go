// Package auth provides authentication and authorization services
package auth

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new auth handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

// MagicLinkRequest represents a magic link request
type MagicLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// OAuthCallbackRequest represents OAuth callback parameters
type OAuthCallbackRequest struct {
	Code  string `query:"code"  validate:"required"`
	State string `query:"state" validate:"required"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// APITokenRequest represents an API token creation request
type APITokenRequest struct {
	Name   string   `json:"name"   validate:"required"`
	Scopes []string `json:"scopes"`
}

// Login handles password-based login
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	session, err := h.service.AuthenticateWithPassword(
		c.Request().Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		h.logger.Error("Login failed", "email", req.Email, "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	})
}

// RequestMagicLink handles magic link requests
func (h *Handler) RequestMagicLink(c echo.Context) error {
	var req MagicLinkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	err := h.service.AuthenticateWithMagicLink(c.Request().Context(), req.Email)
	if err != nil {
		h.logger.Error("Magic link request failed", "email", req.Email, "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to send magic link"},
		)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Magic link sent to your email"})
}

// VerifyMagicLink handles magic link verification
func (h *Handler) VerifyMagicLink(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token required"})
	}

	session, err := h.service.VerifyMagicLink(c.Request().Context(), token)
	if err != nil {
		h.logger.Error("Magic link verification failed", "error", err)
		return c.JSON(
			http.StatusUnauthorized,
			map[string]string{"error": "Invalid or expired token"},
		)
	}

	return c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	})
}

// OAuthLogin initiates OAuth login
func (h *Handler) OAuthLogin(c echo.Context) error {
	provider := c.Param("provider")

	// Get OAuth provider from service
	authURL := h.service.GetOAuthURL(Provider(provider))
	if authURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unsupported provider"})
	}

	// Redirect to OAuth provider
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// OAuthCallback handles OAuth callback
func (h *Handler) OAuthCallback(c echo.Context) error {
	provider := c.Param("provider")

	var req OAuthCallbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	session, err := h.service.AuthenticateWithOAuth(
		c.Request().Context(),
		Provider(provider),
		req.Code,
	)
	if err != nil {
		h.logger.Error("OAuth callback failed", "provider", provider, "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication failed"})
	}

	// For web apps, typically redirect to frontend with tokens
	// For API, return tokens directly
	return c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	})
}

// RefreshToken handles token refresh
func (h *Handler) RefreshToken(c echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	session, err := h.service.RefreshSession(c.Request().Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error("Token refresh failed", "error", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	})
}

// CreateAPIToken creates a new API token
func (h *Handler) CreateAPIToken(c echo.Context) error {
	// Get user from context (would be set by auth middleware)
	userID := getUserIDFromContext(c)
	if userID == nil {
		return c.JSON(
			http.StatusUnauthorized,
			map[string]string{"error": "Authentication required"},
		)
	}

	var req APITokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	token, err := h.service.CreateAPIToken(c.Request().Context(), *userID, req.Name, req.Scopes)
	if err != nil {
		h.logger.Error("API token creation failed", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to create token"},
		)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"token": token,
		"name":  req.Name,
	})
}

// Logout handles logout
func (h *Handler) Logout(c echo.Context) error {
	// Get session from context
	sessionID := getSessionIDFromContext(c)
	if sessionID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No active session"})
	}

	err := h.service.RevokeSession(c.Request().Context(), *sessionID)
	if err != nil {
		h.logger.Error("Logout failed", "error", err)
		// Still return success to client
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// LogoutAll revokes all sessions for the user
func (h *Handler) LogoutAll(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == nil {
		return c.JSON(
			http.StatusUnauthorized,
			map[string]string{"error": "Authentication required"},
		)
	}

	err := h.service.RevokeAllSessions(c.Request().Context(), *userID)
	if err != nil {
		h.logger.Error("Logout all failed", "error", err)
		// Still return success to client
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "All sessions revoked"})
}

// Helper functions to extract auth information from context
func getUserIDFromContext(c echo.Context) *uuid.UUID {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok || userID == uuid.Nil {
		return nil
	}
	return &userID
}

func getSessionIDFromContext(c echo.Context) *uuid.UUID {
	// SessionID is stored in JWT subject claim
	claims, ok := GetJWTClaims(c)
	if !ok || claims.Subject == "" {
		return nil
	}
	// Try to parse subject as UUID (session ID)
	sessionID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil
	}
	return &sessionID
}
