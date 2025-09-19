package sessions

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/archesai/archesai/internal/users"
)

// Handler handles HTTP requests for sessions
type Handler struct {
	service      *Service
	usersService *users.Service
	logger       *slog.Logger
}

// NewHandler creates a new sessions handler
func NewHandler(service *Service, usersService *users.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service:      service,
		usersService: usersService,
		logger:       logger,
	}
}

// ListSessions lists all sessions for the authenticated user
func (h *Handler) ListSessions(ctx echo.Context, _ ListSessionsParams) error {
	// Get user from context
	userID := ctx.Get("userID")
	if userID == nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
	}

	// For now, return empty list
	return ctx.JSON(http.StatusOK, ListSessions200JSONResponse{
		Data: []Session{},
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: 0,
		},
	})
}

// CreateSession creates a new session (login)
func (h *Handler) CreateSession(ctx echo.Context) error {
	var req CreateSessionJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid request body",
		})
	}

	// For now, we'll create a mock user since we don't have user lookup by email
	// In production, you would query the user from the database by email
	userID := uuid.New()

	// For demo purposes, accept any email/password combination
	// In production, you would verify against the database
	if req.Password == "" {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "invalid credentials",
		})
	}

	// Hash the password for future reference (not used here)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	_ = hashedPassword

	// Create session
	sessionID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	if req.RememberMe {
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
	}

	// Note: session is created but not stored for simplicity
	_ = &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate tokens
	accessToken, err := h.service.GenerateAccessToken(userID, sessionID)
	if err != nil {
		h.logger.Error("failed to generate access token", "error", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "failed to create session",
		})
	}

	refreshToken, err := h.service.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		h.logger.Error("failed to generate refresh token", "error", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "failed to create session",
		})
	}

	// Store session in memory or database (simplified for now)
	// In production, this should be stored in database

	// Set session cookie
	isHTTPS := ctx.Request().TLS != nil || ctx.Request().Header.Get("X-Forwarded-Proto") == "https"
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   isHTTPS, // Set based on request protocol
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	}
	ctx.SetCookie(cookie)

	// Return token response
	return ctx.JSON(http.StatusCreated, CreateSession201JSONResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

// GetSession gets a specific session
func (h *Handler) GetSession(ctx echo.Context, id UUID) error {
	// For now, return a mock session
	session := Session{
		ID:             id,
		UserID:         uuid.New(),
		OrganizationID: uuid.New(),
		IPAddress:      ctx.RealIP(),
		UserAgent:      ctx.Request().UserAgent(),
		Token:          "mock-token",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	return ctx.JSON(http.StatusOK, GetSession200JSONResponse{
		Data: session,
	})
}

// UpdateSession updates a session
func (h *Handler) UpdateSession(ctx echo.Context, _ UUID) error {
	// Not implemented yet
	return ctx.JSON(http.StatusNotImplemented, map[string]string{
		"message": "not implemented",
	})
}

// DeleteSession deletes a session (logout)
func (h *Handler) DeleteSession(ctx echo.Context, _ UUID) error {
	// Clear session cookie
	isHTTPS := ctx.Request().TLS != nil || ctx.Request().Header.Get("X-Forwarded-Proto") == "https"
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isHTTPS, // Set based on request protocol
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusNoContent)
}
