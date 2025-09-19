package sessions

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// StrictServer implements StrictServerInterface for sessions
type StrictServer struct {
	service *Service
	logger  *slog.Logger
}

// NewStrictServer creates a new strict server implementation
func NewStrictServer(service *Service, logger *slog.Logger) *StrictServer {
	return &StrictServer{
		service: service,
		logger:  logger,
	}
}

// ListSessions lists all sessions for the authenticated user
func (s *StrictServer) ListSessions(
	_ context.Context,
	_ ListSessionsRequestObject,
) (ListSessionsResponseObject, error) {
	// For now, return empty list
	return ListSessions200JSONResponse{
		Data: []Session{},
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: 0,
		},
	}, nil
}

// CreateSession creates a new session (login)
func (s *StrictServer) CreateSession(
	_ context.Context,
	request CreateSessionRequestObject,
) (CreateSessionResponseObject, error) {
	if request.Body == nil {
		return CreateSession400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Detail: "Request body is required",
				Status: 400,
				Title:  "Bad Request",
				Type:   "about:blank",
			},
		}, nil
	}

	// For demo purposes, accept any email/password combination
	// In production, you would verify against the database
	if request.Body.Password == "" {
		return CreateSession401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
				Detail: "Invalid credentials",
				Status: 401,
				Title:  "Unauthorized",
				Type:   "about:blank",
			},
		}, nil
	}

	// Create mock user for demo
	userID := uuid.New()
	sessionID := uuid.New()

	// Determine session expiration
	if request.Body.RememberMe {
		// Long-lived session
		_ = time.Now().Add(30 * 24 * time.Hour)
	} else {
		// Normal session
		_ = time.Now().Add(24 * time.Hour)
	}

	// Generate tokens
	accessToken, err := s.service.GenerateAccessToken(userID, sessionID)
	if err != nil {
		s.logger.Error("failed to generate access token", "error", err)
		return CreateSession401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
				Detail: "Failed to create session",
				Status: 500,
				Title:  "Internal Server Error",
				Type:   "about:blank",
			},
		}, nil
	}

	refreshToken, err := s.service.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		s.logger.Error("failed to generate refresh token", "error", err)
		return CreateSession401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: UnauthorizedApplicationProblemPlusJSONResponse{
				Detail: "Failed to create session",
				Status: 500,
				Title:  "Internal Server Error",
				Type:   "about:blank",
			},
		}, nil
	}

	// Return success response
	return CreateSession201JSONResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

// GetSession gets a specific session
func (s *StrictServer) GetSession(
	_ context.Context,
	request GetSessionRequestObject,
) (GetSessionResponseObject, error) {
	// For now, return a mock session
	session := Session{
		ID:             request.ID,
		UserID:         uuid.New(),
		OrganizationID: uuid.New(),
		IPAddress:      "127.0.0.1",
		UserAgent:      "Mozilla/5.0",
		Token:          "mock-token",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	return GetSession200JSONResponse{
		Data: session,
	}, nil
}

// UpdateSession updates a session
func (s *StrictServer) UpdateSession(
	_ context.Context,
	request UpdateSessionRequestObject,
) (UpdateSessionResponseObject, error) {
	// Not implemented yet
	session := Session{
		ID:             request.ID,
		UserID:         uuid.New(),
		OrganizationID: uuid.New(),
		IPAddress:      "127.0.0.1",
		UserAgent:      "Mozilla/5.0",
		Token:          "mock-token",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	return UpdateSession200JSONResponse{
		Data: session,
	}, nil
}

// DeleteSession deletes a session (logout)
func (s *StrictServer) DeleteSession(
	_ context.Context,
	request DeleteSessionRequestObject,
) (DeleteSessionResponseObject, error) {
	// For now, return the deleted session
	session := Session{
		ID:             request.ID,
		UserID:         uuid.New(),
		OrganizationID: uuid.New(),
		IPAddress:      "127.0.0.1",
		UserAgent:      "Mozilla/5.0",
		Token:          "deleted",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now(), // Already expired
	}

	return DeleteSession200JSONResponse{
		Data: session,
	}, nil
}
