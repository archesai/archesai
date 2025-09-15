package sessions

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SessionService defines the session service interface
type SessionService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, ipAddress, userAgent string, rememberMe bool) (*Session, string, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	FindSessions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Session, int64, error)
	FindSessionByID(ctx context.Context, sessionID uuid.UUID) (*Session, error)
}

// UserService defines the user authentication interface
type UserService interface {
	AuthenticateUser(ctx context.Context, email, password string) (uuid.UUID, uuid.UUID, error) // returns userID, organizationID, error
}

// Handler handles HTTP requests for sessions
type Handler struct {
	sessionService SessionService
	userService    UserService
	logger         *slog.Logger
}

// NewHandler creates a new session handler
func NewHandler(sessionService SessionService, userService UserService, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		sessionService: sessionService,
		userService:    userService,
		logger:         logger,
	}
}

// SessionsCreate handles POST /auth/sessions (login)
func (h *Handler) SessionsCreate(ctx echo.Context) error {
	var req CreateSessionJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		h.logger.Error("failed to bind login request", "error", err)
		return ctx.JSON(http.StatusBadRequest, Problem{
			Type:   "validation_failed",
			Title:  "Invalid Request",
			Status: http.StatusBadRequest,
			Detail: "Failed to parse request body",
		})
	}

	// Authenticate user with email and password
	userID, organizationID, err := h.userService.AuthenticateUser(ctx.Request().Context(), string(req.Email), req.Password)
	if err != nil {
		h.logger.Warn("authentication failed", "email", req.Email, "error", err)
		return ctx.JSON(http.StatusUnauthorized, Problem{
			Type:   "authentication_failed",
			Title:  "Authentication Failed",
			Status: http.StatusUnauthorized,
			Detail: "Invalid email or password",
		})
	}

	// Create session
	ipAddress := ctx.RealIP()
	userAgent := ctx.Request().Header.Get("User-Agent")
	rememberMe := req.RememberMe

	session, token, err := h.sessionService.CreateSession(
		ctx.Request().Context(),
		userID,
		organizationID,
		ipAddress,
		userAgent,
		rememberMe,
	)
	if err != nil {
		h.logger.Error("failed to create session", "error", err, "userId", userID)
		return ctx.JSON(http.StatusInternalServerError, Problem{
			Type:   "internal_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "Failed to create session",
		})
	}

	// Return token response
	tokenResponse := TokenResponse{
		AccessToken:  token,
		RefreshToken: token, // For simplicity, using same token (in practice, implement refresh tokens)
		TokenType:    "bearer",
		ExpiresIn:    int64(24 * 60 * 60), // 24 hours in seconds
	}

	h.logger.Info("user logged in successfully", "userId", userID, "sessionId", session.Id)
	return ctx.JSON(http.StatusCreated, tokenResponse)
}

// SessionsList handles GET /auth/sessions
func (h *Handler) SessionsList(ctx echo.Context, params ListSessionsParams) error {
	// TODO: Extract user ID from authentication context
	// For now, using a placeholder
	userID := uuid.New()

	limit := 20
	offset := 0

	// Extract pagination from Page parameter
	if params.Page.Size > 0 {
		limit = params.Page.Size
	}
	if params.Page.Number > 0 {
		offset = (params.Page.Number - 1) * limit
	}

	sessions, total, err := h.sessionService.FindSessions(ctx.Request().Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("failed to find sessions", "error", err, "userId", userID)
		return ctx.JSON(http.StatusInternalServerError, Problem{
			Type:   "internal_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "Failed to retrieve sessions",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": sessions,
		"meta": map[string]interface{}{
			"total": total,
		},
	})
}

// SessionsDelete handles DELETE /auth/sessions/{id} (logout)
func (h *Handler) SessionsDelete(ctx echo.Context, id uuid.UUID) error {
	err := h.sessionService.DeleteSession(ctx.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to delete session", "error", err, "sessionId", id)
		return ctx.JSON(http.StatusInternalServerError, Problem{
			Type:   "internal_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "Failed to delete session",
		})
	}

	h.logger.Info("session deleted successfully", "sessionId", id)
	return ctx.NoContent(http.StatusNoContent)
}

// SessionsGetOne handles GET /auth/sessions/{id}
func (h *Handler) SessionsGetOne(ctx echo.Context, id uuid.UUID) error {
	session, err := h.sessionService.FindSessionByID(ctx.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to find session", "error", err, "sessionId", id)
		return ctx.JSON(http.StatusNotFound, Problem{
			Type:   "not_found",
			Title:  "Session Not Found",
			Status: http.StatusNotFound,
			Detail: "The requested session does not exist",
		})
	}

	return ctx.JSON(http.StatusOK, session)
}
