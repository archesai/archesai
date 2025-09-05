// Package http provides HTTP handlers for authentication management.
package http

import (
	"log/slog"
	"net/http"

	"github.com/archesai/archesai/internal/domains/auth/core"
	"github.com/archesai/archesai/internal/domains/auth/generated/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests for auth operations
type Handler struct {
	service *core.Service
	logger  *slog.Logger
}

// NewHandler creates a new auth HTTP handler
func NewHandler(service *core.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Register handles user registration (implements ServerInterface)
func (h *Handler) Register(c echo.Context) error {
	var req api.RegisterJSONBody
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, tokens, err := h.service.SignUp(c.Request().Context(), &core.SignUpRequest{
		Email:    string(req.Email),
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		switch err {
		case core.ErrUserExists:
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

// Login handles user authentication (implements ServerInterface)
func (h *Handler) Login(c echo.Context) error {
	var req api.LoginJSONBody
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	user, tokens, err := h.service.SignIn(c.Request().Context(), &core.SignInRequest{
		Email:    string(req.Email),
		Password: req.Password,
	}, ipAddress, userAgent)
	if err != nil {
		switch err {
		case core.ErrInvalidCredentials:
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
		case core.ErrInvalidToken:
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
		case core.ErrInvalidToken, core.ErrTokenExpired:
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired refresh token")
		case core.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to refresh token", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(http.StatusOK, tokens)
}

// GetOneUser handles getting a single user
func (h *Handler) GetOneUser(ctx echo.Context, id uuid.UUID) error {
	userID := id

	user, err := h.service.GetUser(ctx.Request().Context(), userID)
	if err != nil {
		switch err {
		case core.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to get user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return ctx.JSON(http.StatusOK, user)
}

// UpdateUser handles updating a user
func (h *Handler) UpdateUser(ctx echo.Context, id uuid.UUID) error {
	userID := id

	var req api.UpdateUserJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Map to domain request
	domainReq := &core.UpdateUserRequest{}
	if req.Email != "" {
		// Note: Email update might need special handling (verification, etc.)
		// For now, we'll skip email updates via this endpoint
		h.logger.Info("email update requested but skipped", "user_id", id)
	}
	if req.Image != "" {
		domainReq.Image = &req.Image
	}

	// TODO: Implement UpdateUser with proper type mapping
	// For now, just get the user
	user, err := h.service.GetUser(ctx.Request().Context(), userID)
	if err != nil {
		switch err {
		case core.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to update user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return ctx.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion
func (h *Handler) DeleteUser(ctx echo.Context, id uuid.UUID) error {
	userID := id

	err := h.service.DeleteUser(ctx.Request().Context(), userID)
	if err != nil {
		switch err {
		case core.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		default:
			h.logger.Error("failed to delete user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// FindManyUsers handles listing users with pagination
func (h *Handler) FindManyUsers(ctx echo.Context, params api.FindManyUsersParams) error {
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
	e.POST("/auth/sign-up", h.Register)
	e.POST("/auth/sign-in", h.Login)
	e.POST("/auth/sign-out", h.SignOut)
	e.POST("/auth/refresh", h.RefreshToken)

	// User CRUD routes
	e.GET("/auth/users", func(ctx echo.Context) error {
		// Parse query parameters into FindManyUsersParams
		var params api.FindManyUsersParams
		if err := (&echo.DefaultBinder{}).BindQueryParams(ctx, &params); err != nil {
			return err
		}
		return h.FindManyUsers(ctx, params)
	})
	e.GET("/auth/users/:id", func(ctx echo.Context) error {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
		}
		return h.GetOneUser(ctx, id)
	})
	e.PATCH("/auth/users/:id", func(ctx echo.Context) error {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
		}
		return h.UpdateUser(ctx, id)
	})
	e.DELETE("/auth/users/:id", func(ctx echo.Context) error {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID")
		}
		return h.DeleteUser(ctx, id)
	})
}

// Helper converter functions

// convertToGeneratedUsers converts a slice of domain Users to generated UserEntity
func convertToGeneratedUsers(domainUsers []*core.User) []api.UserEntity {
	result := make([]api.UserEntity, len(domainUsers))
	for i, u := range domainUsers {
		result[i] = u.UserEntity
	}
	return result
}

// convertPagination converts generated pagination params to domain options
func convertPagination(page api.Page) (limit, offset int32) {
	limit = 50 // default
	offset = 0 // default

	if page.Size > 0 {
		limit = int32(page.Size)
		if limit > 100 {
			limit = 100 // max limit
		}
	}
	if page.Number > 0 {
		offset = int32(page.Number-1) * limit
	}

	return limit, offset
}
