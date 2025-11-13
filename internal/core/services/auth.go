// Package services provides the authentication service interface.
package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/core/models"
)

// AuthService defines the authentication service interface.
type AuthService interface {
	// Core authentication methods
	Register(ctx context.Context, email, password, name string) (*models.AuthTokens, error)
	AuthenticateWithPassword(
		ctx context.Context,
		email, password string,
	) (*models.AuthTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthTokens, error)

	// OAuth methods
	GetOAuthAuthorizationURL(provider string, state string) (string, error)
	HandleOAuthCallback(
		ctx context.Context,
		provider string,
		code string,
		state string,
	) (*models.AuthTokens, error)

	// Magic link methods
	GenerateMagicLink(ctx context.Context, identifier, redirectURL string) (string, error)
	VerifyMagicLink(ctx context.Context, token string) (*models.AuthTokens, error)

	// Session management
	GetSessionByToken(ctx context.Context, accessToken string) (*models.Session, error)
	DeleteSessionByID(ctx context.Context, sessionID uuid.UUID) error
	DeleteAllUserSessions(ctx context.Context, sessionID uuid.UUID) error

	// Password reset
	RequestPasswordReset(ctx context.Context, email string) error
	ConfirmPasswordReset(ctx context.Context, token, newPassword string) error

	// Email verification
	RequestEmailVerification(ctx context.Context, sessionID uuid.UUID) error
	ConfirmEmailVerification(ctx context.Context, token string) error

	// Email change
	RequestEmailChange(ctx context.Context, sessionID uuid.UUID, newEmail string) error
	ConfirmEmailChange(ctx context.Context, token, newEmail string, userID uuid.UUID) error

	// Account management
	LinkAccount(
		ctx context.Context,
		sessionID uuid.UUID,
		provider string,
		accessToken *string,
	) (*models.Account, error)
	DeleteAccount(ctx context.Context, sessionID uuid.UUID) (*models.Account, error)
	UpdateAccount(
		ctx context.Context,
		sessionID uuid.UUID,
		updates map[string]any,
	) (*models.User, error)
}
