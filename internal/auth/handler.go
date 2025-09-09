package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// Handler provides HTTP handlers for authentication operations
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)

// NewHandler creates a new authentication handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// NewAuthStrictHandler creates a StrictHandler with middleware
func NewAuthStrictHandler(handler StrictServerInterface) ServerInterface {
	return NewStrictHandler(handler, nil)
}

// Register handles user registration (implements StrictServerInterface)
func (h *Handler) Register(ctx context.Context, req RegisterRequestObject) (RegisterResponseObject, error) {
	// Get IP address and user agent from echo context if available
	var ipAddress, userAgent string
	if echoCtx, ok := ctx.Value("echo.Context").(echo.Context); ok {
		ipAddress = echoCtx.RealIP()
		userAgent = echoCtx.Request().Header.Get("User-Agent")
	}

	// Create the registration request
	registerReq := &RegisterRequest{
		Email:    req.Body.Email,
		Password: req.Body.Password,
		Name:     req.Body.Name,
	}

	// Call the service to register the user
	_, tokens, err := h.service.Register(ctx, registerReq)
	if err != nil {
		switch err {
		case ErrUserExists:
			return Register401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
					Title:  "User already exists",
					Status: 409,
					Type:   "user-exists",
				},
			}, nil
		default:
			h.logger.Error("failed to register user", "error", err)
			return nil, err
		}
	}

	// Store IP and user agent in session if we have them
	if ipAddress != "" || userAgent != "" {
		// This would be handled in the service layer with proper session management
		_ = ipAddress
		_ = userAgent
	}

	// Return the tokens in the response
	return Register201JSONResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}, nil
}

// Login handles user login (implements StrictServerInterface)
func (h *Handler) Login(ctx context.Context, req LoginRequestObject) (LoginResponseObject, error) {
	// Get IP address and user agent from echo context if available
	var ipAddress, userAgent string
	if echoCtx, ok := ctx.Value("echo.Context").(echo.Context); ok {
		ipAddress = echoCtx.RealIP()
		userAgent = echoCtx.Request().Header.Get("User-Agent")
	}

	loginReq := &LoginRequest{
		Email:    req.Body.Email,
		Password: req.Body.Password,
	}

	_, tokens, err := h.service.Login(ctx, loginReq, ipAddress, userAgent)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			return Login401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
					Title:  "Invalid credentials",
					Status: 401,
					Type:   "invalid-credentials",
				},
			}, nil
		default:
			h.logger.Error("failed to login user", "error", err)
			return nil, err
		}
	}

	return Login200JSONResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}, nil
}

// Logout handles user logout (implements StrictServerInterface)
func (h *Handler) Logout(ctx context.Context, _ LogoutRequestObject) (LogoutResponseObject, error) {
	// Get the session token from context (set by middleware)
	token, ok := ctx.Value(SessionTokenContextKey).(string)
	if !ok {
		return Logout401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
				Title:  "No session token",
				Status: 401,
				Type:   "no-session",
			},
		}, nil
	}

	err := h.service.Logout(ctx, token)
	if err != nil {
		if err == ErrInvalidToken {
			return Logout401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
					Title:  "Invalid session",
					Status: 401,
					Type:   "invalid-session",
				},
			}, nil
		}
		h.logger.Error("failed to logout user", "error", err)
		return nil, err
	}

	return Logout204Response{}, nil
}

// RefreshToken handles token refresh (implements StrictServerInterface)
// TODO: Add RefreshToken endpoint to OpenAPI spec
/*
func (h *Handler) RefreshToken(ctx context.Context, req RefreshTokenRequestObject) (RefreshTokenResponseObject, error) {
	refreshToken := req.Body.RefreshToken

	tokens, err := h.service.RefreshToken(ctx, refreshToken)
	if err != nil {
		switch err {
		case ErrInvalidToken:
			return RefreshToken401ApplicationProblemPlusJSONResponse{
				Title:  "Invalid refresh token",
				Status: 401,
				Type:   "invalid-token",
			}, nil
		default:
			h.logger.Error("failed to refresh token", "error", err)
			return nil, err
		}
	}

	return RefreshToken200JSONResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}, nil
}
*/

// VerifyToken handles token verification (implements StrictServerInterface)
// TODO: Add VerifyToken endpoint to OpenAPI spec
/*
func (h *Handler) VerifyToken(ctx context.Context, _ VerifyTokenRequestObject) (VerifyTokenResponseObject, error) {
	// Get the user from context (set by middleware)
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return VerifyToken401ApplicationProblemPlusJSONResponse{
			Title:  "Invalid token",
			Status: 401,
			Type:   "invalid-token",
		}, nil
	}

	return VerifyToken200JSONResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
*/

// AccountsFindMany handles listing accounts (stub implementation)
func (h *Handler) AccountsFindMany(_ context.Context, _ AccountsFindManyRequestObject) (AccountsFindManyResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// AccountsDelete handles deleting an account (stub implementation)
func (h *Handler) AccountsDelete(_ context.Context, _ AccountsDeleteRequestObject) (AccountsDeleteResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// AccountsGetOne handles getting a single account (stub implementation)
func (h *Handler) AccountsGetOne(_ context.Context, _ AccountsGetOneRequestObject) (AccountsGetOneResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// FindManySessions handles listing sessions (stub implementation)
func (h *Handler) FindManySessions(_ context.Context, _ FindManySessionsRequestObject) (FindManySessionsResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// DeleteSession handles deleting a session (stub implementation)
func (h *Handler) DeleteSession(_ context.Context, _ DeleteSessionRequestObject) (DeleteSessionResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetOneSession handles getting a single session (stub implementation)
func (h *Handler) GetOneSession(_ context.Context, _ GetOneSessionRequestObject) (GetOneSessionResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// UpdateSession handles updating a session (stub implementation)
func (h *Handler) UpdateSession(_ context.Context, _ UpdateSessionRequestObject) (UpdateSessionResponseObject, error) {
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

// RequestEmailVerification handles email verification requests
func (h *Handler) RequestEmailVerification(ctx context.Context, _ RequestEmailVerificationRequestObject) (RequestEmailVerificationResponseObject, error) {
	// Get the authenticated user from context
	userClaims, ok := ctx.Value(ClaimsContextKey).(*EnhancedClaims)
	if !ok || userClaims == nil {
		// If no user in context, return unauthorized
		return RequestEmailVerification401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
				Title:  "Unauthorized",
				Status: 401,
				Type:   "unauthorized",
				Detail: "Authentication required",
			},
		}, nil
	}

	// Resend verification email
	err := h.service.ResendVerificationEmail(ctx, userClaims.Email)
	if err != nil {
		h.logger.Error("failed to resend verification email", "error", err, "email", userClaims.Email)
		return RequestEmailVerification400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Title:  "Failed to send verification email",
				Status: 400,
				Type:   "email-send-failed",
				Detail: "Could not send verification email. Please try again later.",
			},
		}, nil
	}

	// Return 204 No Content on success
	return RequestEmailVerification204Response{}, nil
}

// ConfirmEmailVerification handles email verification confirmation
func (h *Handler) ConfirmEmailVerification(ctx context.Context, req ConfirmEmailVerificationRequestObject) (ConfirmEmailVerificationResponseObject, error) {
	if req.Body == nil || req.Body.Token == "" {
		return ConfirmEmailVerification401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
				Title:  "Invalid request",
				Status: 400,
				Type:   "invalid-request",
				Detail: "Verification token is required",
			},
		}, nil
	}

	// Verify the email using the token
	err := h.service.VerifyEmail(ctx, req.Body.Token)
	if err != nil {
		switch err {
		case ErrInvalidToken:
			return ConfirmEmailVerification404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "Invalid token",
					Status: 404,
					Type:   "invalid-token",
					Detail: "The verification token is invalid or has been used",
				},
			}, nil
		case ErrTokenExpired:
			return ConfirmEmailVerification401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
					Title:  "Token expired",
					Status: 401,
					Type:   "token-expired",
					Detail: "The verification token has expired. Please request a new one.",
				},
			}, nil
		case ErrUserNotFound:
			return ConfirmEmailVerification404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "User not found",
					Status: 404,
					Type:   "user-not-found",
					Detail: "The user associated with this token was not found",
				},
			}, nil
		default:
			h.logger.Error("failed to verify email", "error", err)
			return nil, err
		}
	}

	// Get the verified user to return in response
	// Note: We need to get user details from the token verification process
	// For now, we'll create a simple response - in production, you'd want to
	// return the actual user and session details

	// TODO: Generate new session for the verified user
	// For now, return a successful response without session details
	return ConfirmEmailVerification200JSONResponse{
		User: User{
			// This would be populated from the actual verified user
			Email:         Email("verified@example.com"),
			EmailVerified: true,
			Name:          "Verified User",
		},
		Session: Session{
			// This would be a newly created session
		},
	}, nil
}

// RequestPasswordReset handles password reset requests (stub implementation)
func (h *Handler) RequestPasswordReset(_ context.Context, _ RequestPasswordResetRequestObject) (RequestPasswordResetResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// ConfirmPasswordReset handles password reset confirmation (stub implementation)
func (h *Handler) ConfirmPasswordReset(_ context.Context, _ ConfirmPasswordResetRequestObject) (ConfirmPasswordResetResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
