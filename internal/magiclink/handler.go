// Package magiclink provides magic link authentication functionality
package magiclink

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/sessions"
	"github.com/archesai/archesai/internal/users"
)

// Handler handles HTTP requests for magic link authentication
type Handler struct {
	service         *Service
	sessionsService *sessions.Service
	usersService    *users.Service
	logger          *slog.Logger
}

// NewHandler creates a new magic link handler
func NewHandler(
	service *Service,
	sessionsService *sessions.Service,
	usersService *users.Service,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		service:         service,
		sessionsService: sessionsService,
		usersService:    usersService,
		logger:          logger,
	}
}

// RequestMagicLink handles magic link generation requests
func (h *Handler) RequestMagicLink(ctx echo.Context) error {
	var req struct {
		Identifier     string `json:"identifier"`
		DeliveryMethod string `json:"deliveryMethod"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate identifier (should be an email)
	if req.Identifier == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Identifier is required",
		})
	}

	// Default to console delivery if not specified
	if req.DeliveryMethod == "" {
		req.DeliveryMethod = "console"
	}

	// Check if user exists (optional - can create on first login)
	user, err := h.usersService.GetByEmail(context.Background(), req.Identifier)
	var userID *uuid.UUID
	if err == nil && user != nil {
		userID = &user.ID
	}

	// Get IP and user agent
	ipAddress := ctx.RealIP()
	userAgent := ctx.Request().UserAgent()

	// Request the magic link
	token, err := h.service.RequestMagicLink(
		context.Background(),
		req.Identifier,
		DeliveryMethod(req.DeliveryMethod),
		userID,
		ipAddress,
		userAgent,
	)

	if err != nil {
		if err == ErrRateLimitExceeded {
			return ctx.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests. Please try again later.",
			})
		}
		h.logger.Error("Failed to generate magic link", "error", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate magic link",
		})
	}

	// Return success response
	response := map[string]interface{}{
		"message": "Magic link sent successfully",
		"method":  req.DeliveryMethod,
	}

	// Include OTP code in response for OTP method
	if req.DeliveryMethod == "otp" && token.Code != "" {
		response["code"] = token.Code // In production, don't return this
	}

	return ctx.JSON(http.StatusOK, response)
}

// VerifyMagicLink handles magic link verification
func (h *Handler) VerifyMagicLink(ctx echo.Context) error {
	var req struct {
		Token      string `json:"token"`
		Code       string `json:"code"`
		Identifier string `json:"identifier"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	var token *Token
	var err error

	// Verify by token or OTP
	if req.Token != "" {
		token, err = h.service.VerifyToken(context.Background(), req.Token)
	} else if req.Code != "" && req.Identifier != "" {
		token, err = h.service.VerifyOTP(context.Background(), req.Identifier, req.Code)
	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Either token or code+identifier is required",
		})
	}

	if err != nil {
		switch err {
		case ErrTokenNotFound, ErrInvalidOTP:
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired magic link",
			})
		case ErrTokenExpired:
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Magic link has expired",
			})
		case ErrTokenAlreadyUsed:
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Magic link has already been used",
			})
		default:
			h.logger.Error("Failed to verify magic link", "error", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to verify magic link",
			})
		}
	}

	// Get or create user
	var user *users.User
	if token.UserID != nil {
		user, err = h.usersService.GetByID(context.Background(), *token.UserID)
		if err != nil {
			h.logger.Error("Failed to get user", "error", err, "userID", token.UserID)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get user",
			})
		}
	} else {
		// Create new user if doesn't exist
		user, err = h.usersService.GetByEmail(context.Background(), token.Identifier)
		if err != nil {
			// User doesn't exist, create new one
			createReq := &users.CreateUserRequest{
				Email:         token.Identifier,
				Name:          strings.Split(token.Identifier, "@")[0], // Use email prefix as name
				EmailVerified: true,                                    // Mark as verified since they used magic link
			}

			user, err = h.usersService.Create(context.Background(), createReq)
			if err != nil {
				h.logger.Error("Failed to create user", "error", err)
				return ctx.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to create user",
				})
			}
		}
	}

	// Create session
	session, err := h.sessionsService.CreateSessionWithMethod(
		context.Background(),
		user.ID,
		nil, // No organization ID initially
		ctx.RealIP(),
		ctx.Request().UserAgent(),
		"magic_link",
		"",
	)
	if err != nil {
		h.logger.Error("Failed to create session", "error", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create session",
		})
	}

	// Generate JWT tokens
	accessToken, err := h.sessionsService.GenerateAccessToken(user.ID, session.ID)
	if err != nil {
		h.logger.Error("Failed to generate access token", "error", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate tokens",
		})
	}

	refreshToken, err := h.sessionsService.GenerateRefreshToken(user.ID, session.ID)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", "error", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate tokens",
		})
	}

	// Set session cookie - use session ID for consistency
	isHTTPS := ctx.Request().TLS != nil || ctx.Request().Header.Get("X-Forwarded-Proto") == "https"
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   isHTTPS, // Set based on request protocol
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 30, // 30 days
	}
	ctx.SetCookie(cookie)

	// Return tokens
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"tokenType":    "Bearer",
		"expiresIn":    3600,
		"user": map[string]interface{}{
			"id":            user.ID,
			"email":         user.Email,
			"name":          user.Name,
			"emailVerified": user.EmailVerified,
		},
	})
}
