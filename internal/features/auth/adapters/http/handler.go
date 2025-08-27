package http

import (
	"net/http"
	"strconv"

	"github.com/archesai/archesai/internal/features/auth/domain"
	"github.com/archesai/archesai/internal/features/auth/ports"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for auth operations
type Handler struct {
	service ports.Service
	logger  *zap.Logger
}

// NewHandler creates a new auth HTTP handler
func NewHandler(service ports.Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// SignUp handles user registration
func (h *Handler) SignUp(c echo.Context) error {
	var req domain.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, tokens, err := h.service.SignUp(c.Request().Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrUserExists:
			return echo.NewHTTPError(http.StatusConflict, "User already exists")
		default:
			h.logger.Error("failed to sign up user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	response := map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	}

	return c.JSON(http.StatusCreated, response)
}

// SignIn handles user authentication
func (h *Handler) SignIn(c echo.Context) error {
	var req domain.SignInRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	user, tokens, err := h.service.SignIn(c.Request().Context(), &req, ipAddress, userAgent)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		default:
			h.logger.Error("failed to sign in user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	response := map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	}

	return c.JSON(http.StatusOK, response)
}

// SignOut handles user logout
func (h *Handler) SignOut(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing authorization token")
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := h.service.SignOut(c.Request().Context(), token)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken:
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		default:
			h.logger.Error("failed to sign out user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Signed out successfully"})
}

// RefreshToken handles token refresh
func (h *Handler) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tokens, err := h.service.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken, domain.ErrTokenExpired:
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired refresh token")
		case domain.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to refresh token", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(http.StatusOK, tokens)
}

// GetUser handles retrieving user information
func (h *Handler) GetUser(c echo.Context) error {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.service.GetUser(c.Request().Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to get user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser handles updating user information
func (h *Handler) UpdateUser(c echo.Context) error {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var req domain.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	user, err := h.service.UpdateUser(c.Request().Context(), userID, &req)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to update user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion
func (h *Handler) DeleteUser(c echo.Context) error {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	err = h.service.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to delete user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// ListUsers handles listing users with pagination
func (h *Handler) ListUsers(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")

	var limit, offset int32 = 50, 0

	if limitParam != "" {
		if l, err := strconv.ParseInt(limitParam, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	if offsetParam != "" {
		if o, err := strconv.ParseInt(offsetParam, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	users, err := h.service.ListUsers(c.Request().Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	response := map[string]interface{}{
		"users":  users,
		"limit":  limit,
		"offset": offset,
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers auth routes with the Echo router
func (h *Handler) RegisterRoutes(e *echo.Group) {
	// Authentication routes
	e.POST("/auth/signup", h.SignUp)
	e.POST("/auth/signin", h.SignIn)
	e.POST("/auth/signout", h.SignOut)
	e.POST("/auth/refresh", h.RefreshToken)

	// User management routes
	e.GET("/users/:id", h.GetUser)
	e.PUT("/users/:id", h.UpdateUser)
	e.DELETE("/users/:id", h.DeleteUser)
	e.GET("/users", h.ListUsers)
}