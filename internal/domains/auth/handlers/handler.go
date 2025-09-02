package handlers

import (
	"log/slog"
	"net/http"

	"github.com/archesai/archesai/internal/domains/auth/entities"
	"github.com/archesai/archesai/internal/domains/auth/services"
	"github.com/archesai/archesai/internal/generated/api/auth/users"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Handler handles HTTP requests for auth operations
// Implements users.ServerInterface
type Handler struct {
	service *services.Service
	logger  *slog.Logger
}

// Ensure Handler implements users.ServerInterface
var _ users.ServerInterface = (*Handler)(nil)

// NewHandler creates a new auth HTTP handler
func NewHandler(service *services.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// SignUp handles user registration
func (h *Handler) SignUp(c echo.Context) error {
	var req entities.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, tokens, err := h.service.SignUp(c.Request().Context(), &req)
	if err != nil {
		switch err {
		case entities.ErrUserExists:
			return echo.NewHTTPError(http.StatusConflict, "User already exists")
		default:
			h.logger.Error("failed to sign up user", "error", err)
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
	var req entities.SignInRequest
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
		case entities.ErrInvalidCredentials:
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		default:
			h.logger.Error("failed to sign in user", "error", err)
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
		case entities.ErrInvalidToken:
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		default:
			h.logger.Error("failed to sign out user", "error", err)
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
		case entities.ErrInvalidToken, entities.ErrTokenExpired:
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired refresh token")
		case entities.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to refresh token", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(http.StatusOK, tokens)
}

// GetOneUser handles retrieving user information
// Implements users.ServerInterface.GetOneUser
func (h *Handler) GetOneUser(ctx echo.Context, id openapi_types.UUID) error {
	userID := uuid.UUID(id)

	user, err := h.service.GetUser(ctx.Request().Context(), userID)
	if err != nil {
		switch err {
		case entities.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to get user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return ctx.JSON(http.StatusOK, user)
}

// UpdateUser handles updating user information
// Implements users.ServerInterface.UpdateUser
func (h *Handler) UpdateUser(ctx echo.Context, id openapi_types.UUID) error {
	userID := uuid.UUID(id)

	var req users.UpdateUserJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Map to domain request
	domainReq := &entities.UpdateUserRequest{}
	if req.Email != nil {
		// Note: Email update might need special handling (verification, etc.)
		// For now, we'll skip email updates via this endpoint
	}
	if req.Image != nil {
		domainReq.Image = req.Image
	}

	user, err := h.service.UpdateUser(ctx.Request().Context(), userID, domainReq)
	if err != nil {
		switch err {
		case entities.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to update user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return ctx.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion
// Implements users.ServerInterface.DeleteUser
func (h *Handler) DeleteUser(ctx echo.Context, id openapi_types.UUID) error {
	userID := uuid.UUID(id)

	err := h.service.DeleteUser(ctx.Request().Context(), userID)
	if err != nil {
		switch err {
		case entities.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to delete user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// FindManyUsers handles listing users with pagination
// Implements users.ServerInterface.FindManyUsers
func (h *Handler) FindManyUsers(ctx echo.Context, params users.FindManyUsersParams) error {
	// Use converter functions for pagination
	limit, offset := convertPagination(params.Page)

	// TODO: Apply filter and sort if needed
	// filter := convertFilter(params.Filter)
	// orderBy, orderDir := convertSort(params.Sort)

	domainUsers, err := h.service.ListUsers(ctx.Request().Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// Convert domain users to generated types
	userEntities := convertToGeneratedUsers(domainUsers)

	response := map[string]interface{}{
		"data":   userEntities,
		"total":  len(userEntities),
		"limit":  limit,
		"offset": offset,
	}

	return ctx.JSON(http.StatusOK, response)
}

// RegisterRoutes registers auth routes with the Echo router
func (h *Handler) RegisterRoutes(e *echo.Group) {
	// Authentication routes
	e.POST("/auth/signup", h.SignUp)
	e.POST("/auth/signin", h.SignIn)
	e.POST("/auth/signout", h.SignOut)
	e.POST("/auth/refresh", h.RefreshToken)

	// Note: User routes will be registered via the generated RegisterHandlers function
}
