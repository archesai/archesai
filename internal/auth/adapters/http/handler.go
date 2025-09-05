// Package http provides HTTP handlers for authentication operations
package http

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/archesai/archesai/internal/auth/domain"
	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

// Context keys for request metadata
const (
	ipAddressKey contextKey = "ip_address"
	userAgentKey contextKey = "user_agent"
	authTokenKey contextKey = "auth_token"
)

// Handler handles HTTP requests for auth operations
type Handler struct {
	service *domain.Service
	logger  *slog.Logger
}

// NewHandler creates a new auth HTTP handler
func NewHandler(service *domain.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Register handles user registration (implements StrictServerInterface)
func (h *Handler) Register(ctx context.Context, req RegisterRequestObject) (RegisterResponseObject, error) {
	user, _, err := h.service.Register(ctx, &domain.RegisterRequest{
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
		Name:     req.Body.Name,
	})
	if err != nil {
		switch err {
		case domain.ErrUserExists:
			// Return 401 Unauthorized (there's no 409 response defined)
			return Register401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
					Detail: "User already exists",
					Status: 401,
					Title:  "User already exists",
				},
			}, nil
		default:
			h.logger.Error("failed to register user", "error", err)
			// Return error for 500 Internal Server Error
			return nil, err
		}
	}

	return Register201JSONResponse{
		Data: user.UserEntity,
	}, nil
}

// Login handles user authentication (implements StrictServerInterface)
func (h *Handler) Login(ctx context.Context, req LoginRequestObject) (LoginResponseObject, error) {
	// Extract IP address and user agent from context (set by middleware)
	ipAddress := "unknown"
	userAgent := "unknown"

	if ip := ctx.Value(ipAddressKey); ip != nil {
		if ipStr, ok := ip.(string); ok {
			ipAddress = ipStr
		}
	}

	if ua := ctx.Value(userAgentKey); ua != nil {
		if uaStr, ok := ua.(string); ok {
			userAgent = uaStr
		}
	}

	user, _, err := h.service.Login(ctx, &domain.LoginRequest{
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
	}, ipAddress, userAgent)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			return Login401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
					Detail: "Invalid credentials",
					Status: 401,
					Title:  "Invalid credentials",
				},
			}, nil
		default:
			h.logger.Error("failed to login user", "error", err)
			return nil, err
		}
	}

	return Login200JSONResponse{
		Data: user.UserEntity,
	}, nil
}

// Logout handles user logout (implements StrictServerInterface)
func (h *Handler) Logout(ctx context.Context, _ LogoutRequestObject) (LogoutResponseObject, error) {
	// Extract token from context (set by auth middleware)
	token := ""
	if t := ctx.Value(authTokenKey); t != nil {
		if tokenStr, ok := t.(string); ok {
			token = tokenStr
		}
	}

	if token == "" {
		return Logout401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
				Detail: "Missing authorization token",
				Status: 401,
				Title:  "Missing authorization token",
			},
		}, nil
	}

	err := h.service.Logout(ctx, token)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken:
			return Logout401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
					Detail: "Invalid token",
					Status: 401,
					Title:  "Invalid token",
				},
			}, nil
		default:
			h.logger.Error("failed to logout user", "error", err)
			return nil, err
		}
	}

	return Logout204Response{}, nil
}

// contextMiddleware injects HTTP request details into context for StrictServerInterface methods
func contextMiddleware(f strictecho.StrictEchoHandlerFunc, _ string) strictecho.StrictEchoHandlerFunc {
	return func(ctx echo.Context, request interface{}) (interface{}, error) {
		// Create new context with request details
		newCtx := context.WithValue(ctx.Request().Context(), ipAddressKey, ctx.RealIP())
		newCtx = context.WithValue(newCtx, userAgentKey, ctx.Request().UserAgent())

		// Extract auth token if present
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader != "" {
			// Remove "Bearer " prefix if present
			token := authHeader
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
			newCtx = context.WithValue(newCtx, authTokenKey, token)
		}

		// Create new request with enriched context
		enrichedCtx := ctx.Request().WithContext(newCtx)
		ctx.SetRequest(enrichedCtx)

		return f(ctx, request)
	}
}

// NewStrictHandlerWithMiddleware creates a StrictHandler with auth-specific middleware
func NewStrictHandlerWithMiddleware(handler StrictServerInterface) ServerInterface {
	return NewStrictHandler(handler, []StrictMiddlewareFunc{contextMiddleware})
}

// GetOneUser handles getting a single user (implements StrictServerInterface)
func (h *Handler) GetOneUser(ctx context.Context, req GetOneUserRequestObject) (GetOneUserResponseObject, error) {
	user, err := h.service.GetUser(ctx, req.Id)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return GetOneUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "User not found",
					Status: 404,
					Type:   "user-not-found",
				},
			}, nil
		default:
			h.logger.Error("failed to get user", "error", err)
			return nil, err
		}
	}

	return GetOneUser200JSONResponse{
		Data: user.UserEntity,
	}, nil
}

// UpdateUser handles updating a user (implements StrictServerInterface)
func (h *Handler) UpdateUser(ctx context.Context, req UpdateUserRequestObject) (UpdateUserResponseObject, error) {
	// TODO: Implement actual user update logic
	// For now, just return the existing user
	user, err := h.service.GetUser(ctx, req.Id)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return UpdateUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "User not found",
					Status: 404,
					Type:   "user-not-found",
				},
			}, nil
		default:
			h.logger.Error("failed to update user", "error", err)
			return nil, err
		}
	}

	return UpdateUser200JSONResponse{
		Data: user.UserEntity,
	}, nil
}

// DeleteUser handles user deletion (implements StrictServerInterface)
func (h *Handler) DeleteUser(ctx context.Context, req DeleteUserRequestObject) (DeleteUserResponseObject, error) {
	err := h.service.DeleteUser(ctx, req.Id)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return DeleteUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "User not found",
					Status: 404,
					Type:   "user-not-found",
				},
			}, nil
		default:
			h.logger.Error("failed to delete user", "error", err)
			return nil, err
		}
	}

	return DeleteUser200JSONResponse{
		Data: domain.UserEntity{Id: req.Id}, // Placeholder response
	}, nil
}

// FindManyUsers handles listing users with pagination (implements StrictServerInterface)
func (h *Handler) FindManyUsers(ctx context.Context, req FindManyUsersRequestObject) (FindManyUsersResponseObject, error) {
	// Use converter functions for pagination
	limit, offset := convertPagination(req.Params.Page)

	domainUsers, err := h.service.ListUsers(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		return nil, err
	}

	// Convert domain users to generated types
	userEntities := convertToGeneratedUsers(domainUsers)

	return FindManyUsers200JSONResponse{
		Data: userEntities,
		Meta: struct {
			Total float32 `json:"total"`
		}{Total: float32(len(userEntities))},
	}, nil
}

// AccountsFindMany handles listing accounts (stub implementation)
func (h *Handler) AccountsFindMany(_ context.Context, _ AccountsFindManyRequestObject) (AccountsFindManyResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// AccountsDelete handles account deletion (stub implementation)
func (h *Handler) AccountsDelete(_ context.Context, _ AccountsDeleteRequestObject) (AccountsDeleteResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// AccountsGetOne handles getting a single account (stub implementation)
func (h *Handler) AccountsGetOne(_ context.Context, _ AccountsGetOneRequestObject) (AccountsGetOneResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// RequestEmailChange handles email change requests (stub implementation)
func (h *Handler) RequestEmailChange(_ context.Context, _ RequestEmailChangeRequestObject) (RequestEmailChangeResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// ConfirmEmailChange handles email change confirmation (stub implementation)
func (h *Handler) ConfirmEmailChange(_ context.Context, _ ConfirmEmailChangeRequestObject) (ConfirmEmailChangeResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// RequestEmailVerification handles email verification requests (stub implementation)
func (h *Handler) RequestEmailVerification(_ context.Context, _ RequestEmailVerificationRequestObject) (RequestEmailVerificationResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// ConfirmEmailVerification handles email verification confirmation (stub implementation)
func (h *Handler) ConfirmEmailVerification(_ context.Context, _ ConfirmEmailVerificationRequestObject) (ConfirmEmailVerificationResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// RequestPasswordReset handles password reset requests (stub implementation)
func (h *Handler) RequestPasswordReset(_ context.Context, _ RequestPasswordResetRequestObject) (RequestPasswordResetResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// ConfirmPasswordReset handles password reset confirmation (stub implementation)
func (h *Handler) ConfirmPasswordReset(_ context.Context, _ ConfirmPasswordResetRequestObject) (ConfirmPasswordResetResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// FindManySessions handles listing sessions (stub implementation)
func (h *Handler) FindManySessions(_ context.Context, _ FindManySessionsRequestObject) (FindManySessionsResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// DeleteSession handles session deletion (stub implementation)
func (h *Handler) DeleteSession(_ context.Context, _ DeleteSessionRequestObject) (DeleteSessionResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetOneSession handles getting a single session (stub implementation)
func (h *Handler) GetOneSession(_ context.Context, _ GetOneSessionRequestObject) (GetOneSessionResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// UpdateSession handles session updates (stub implementation)
func (h *Handler) UpdateSession(_ context.Context, _ UpdateSessionRequestObject) (UpdateSessionResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// Helper converter functions

// convertToGeneratedUsers converts a slice of domain Users to generated UserEntity
func convertToGeneratedUsers(domainUsers []*domain.User) []domain.UserEntity {
	result := make([]domain.UserEntity, len(domainUsers))
	for i, u := range domainUsers {
		result[i] = u.UserEntity
	}
	return result
}

// convertPagination converts generated pagination params to domain options
func convertPagination(page domain.Page) (limit, offset int32) {
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
