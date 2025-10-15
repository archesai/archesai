package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/core/entities"
	corerrors "github.com/archesai/archesai/internal/core/errors"
	"github.com/archesai/archesai/internal/core/repositories"
	"github.com/archesai/archesai/internal/core/services"
	"github.com/archesai/archesai/internal/core/valueobjects"
	"github.com/archesai/archesai/internal/infrastructure/auth/oauth"
	"github.com/archesai/archesai/internal/infrastructure/cache"
	"github.com/archesai/archesai/internal/infrastructure/config"
)

const (
	bindHost = "0.0.0.0"
)

// Service provides authentication functionality across all transport layers.
// It implements the core services.AuthService interface.
type Service struct {
	config             *config.Config
	sessionRepo        repositories.SessionRepository
	sessionsRepo       repositories.SessionRepository
	userRepo           repositories.UserRepository
	accountRepo        repositories.AccountRepository
	accountsRepo       repositories.AccountRepository
	tokenManager       *TokenManager
	magicLink          *MagicLinkProvider
	oauthProviders     map[string]OAuthProvider
	cache              cache.Cache[string]
	jwtSecret          string
	magicLinkDeliverer MagicLinkDeliverer
	oTPDeliverer       OTPDeliverer
}

// MagicLinkDeliverer handles magic link notification delivery.
type MagicLinkDeliverer interface {
	Deliver(ctx context.Context, token *valueobjects.MagicLinkToken, baseURL string) error
}

// OTPDeliverer handles OTP notification delivery.
type OTPDeliverer interface {
	Deliver(ctx context.Context, token *valueobjects.MagicLinkToken, baseURL string) error
}

// Ensure Service implements services.AuthService
var _ services.AuthService = (*Service)(nil)

// OAuthProvider interface for OAuth providers
type OAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth.Tokens, error)
	GetUserInfo(ctx context.Context, accessToken string) (*oauth.UserInfo, error)
}

// NewService creates a new authentication service.
func NewService(
	cfg *config.Config,
	sessionRepo repositories.SessionRepository,
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
	cacheService cache.Cache[string],
	magicLinkDeliverer MagicLinkDeliverer,
	otpDeliverer OTPDeliverer,
) *Service {
	// Build base URL from API config
	baseURL := fmt.Sprintf("http://%s:%d", cfg.API.Host, int(cfg.API.Port))
	if cfg.API.Host == bindHost {
		baseURL = fmt.Sprintf("http://localhost:%d", int(cfg.API.Port))
	}

	s := &Service{
		config:             cfg,
		sessionRepo:        sessionRepo,
		sessionsRepo:       sessionRepo,
		userRepo:           userRepo,
		accountRepo:        accountRepo,
		accountsRepo:       accountRepo,
		tokenManager:       NewTokenManager(cfg.Auth.Local.JWTSecret),
		magicLink:          NewMagicLinkProvider(cfg.Auth.Local.JWTSecret, baseURL),
		oauthProviders:     make(map[string]OAuthProvider),
		cache:              cacheService,
		jwtSecret:          cfg.Auth.Local.JWTSecret,
		magicLinkDeliverer: magicLinkDeliverer,
		oTPDeliverer:       otpDeliverer,
	}

	// Initialize OAuth providers based on config
	if cfg.Auth.Google != nil && cfg.Auth.Google.Enabled && cfg.Auth.Google.ClientID != nil {
		s.oauthProviders["google"] = oauth.NewGoogleProvider(
			*cfg.Auth.Google.ClientID,
			*cfg.Auth.Google.ClientSecret,
			*cfg.Auth.Google.RedirectURL,
		)
	}

	if cfg.Auth.Github != nil && cfg.Auth.Github.Enabled && cfg.Auth.Github.ClientID != nil {
		s.oauthProviders["github"] = oauth.NewGitHubProvider(
			*cfg.Auth.Github.ClientID,
			*cfg.Auth.Github.ClientSecret,
			*cfg.Auth.Github.RedirectURL,
		)
	}

	if cfg.Auth.Microsoft != nil && cfg.Auth.Microsoft.Enabled &&
		cfg.Auth.Microsoft.ClientID != nil {
		s.oauthProviders["microsoft"] = oauth.NewMicrosoftProvider(
			cfg.Auth.Microsoft.ClientID.String(),
			*cfg.Auth.Microsoft.ClientSecret,
			*cfg.Auth.Microsoft.RedirectURL,
		)
	}

	return s
}

// GetOAuthAuthorizationURL generates an OAuth authorization URL.
func (s *Service) GetOAuthAuthorizationURL(provider string, state string) (string, error) {
	p, exists := s.oauthProviders[provider]
	if !exists {
		return "", fmt.Errorf("OAuth provider %s not configured", provider)
	}
	return p.GetAuthURL(state), nil
}

// HandleOAuthCallback processes OAuth callback and creates session.
func (s *Service) HandleOAuthCallback(
	ctx context.Context,
	provider string,
	code string,
	_ string,
) (*valueobjects.AuthTokens, error) {
	p, exists := s.oauthProviders[provider]
	if !exists {
		return nil, fmt.Errorf("OAuth provider %s not configured", provider)
	}

	// Exchange code for tokens
	tokens, err := p.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange OAuth code: %w", err)
	}

	// Get user info from provider
	userInfo, err := p.GetUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Find or create user
	user, err := s.findOrCreateUser(ctx, userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Create session
	session, err := s.createSession(ctx, user.ID, map[string]any{
		"provider":           provider,
		"account_identifier": userInfo.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate auth tokens
	return s.generateTokens(user, session)
}

// GenerateMagicLink creates a stateless magic link token and sends it via the deliverer.
func (s *Service) GenerateMagicLink(
	ctx context.Context,
	identifier, redirectURL string,
) (string, error) {
	link, err := s.magicLink.GenerateLink(identifier, redirectURL)
	if err != nil {
		return "", err
	}

	// Extract token from link for notification
	tokenStr := link[strings.LastIndex(link, "token=")+6:]
	token := &valueobjects.MagicLinkToken{
		Token:      &tokenStr,
		Identifier: identifier,
		ExpiresAt:  time.Now().Add(15 * time.Minute),
	}

	// Send notification if deliverer is configured
	if s.magicLinkDeliverer != nil {
		baseURL := s.config.Platform.URL
		if baseURL == nil || *baseURL == "" {
			return "", fmt.Errorf("platform URL not configured for magic link delivery")
		}
		if err := s.magicLinkDeliverer.Deliver(ctx, token, *baseURL); err != nil {
			return "", fmt.Errorf("failed to deliver magic link: %w", err)
		}
	}

	return link, nil
}

// VerifyMagicLink validates a magic link token and creates a session.
func (s *Service) VerifyMagicLink(
	ctx context.Context,
	token string,
) (*valueobjects.AuthTokens, error) {
	// Validate the magic link token
	claims, err := s.magicLink.ValidateLink(token)
	if err != nil {
		return nil, fmt.Errorf("invalid magic link: %w", err)
	}

	// Find or create user by identifier (email)
	user, err := s.userRepo.GetUserByEmail(ctx, claims.Identifier)
	if err != nil {
		// Only create new user if not found, otherwise propagate the error
		if !errors.Is(err, corerrors.ErrUserNotFound) {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		// Create new user if not exists
		newUser, createErr := entities.NewUser(
			claims.Identifier,
			true,              // Magic link verifies email
			nil,               // No image initially
			claims.Identifier, // Default name to email
		)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create user entity: %w", createErr)
		}

		user, err = s.userRepo.Create(ctx, newUser)
		if err != nil {
			return nil, fmt.Errorf("failed to create user in database: %w", err)
		}

		// Verify the user was actually created by trying to fetch it again
		_, verifyErr := s.userRepo.Get(ctx, user.ID)
		if verifyErr != nil {
			return nil, fmt.Errorf("user creation verification failed: %w", verifyErr)
		}

	}

	session, err := s.createSession(ctx, user.ID, map[string]any{
		"auth_method": "magic_link",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate auth tokens
	return s.generateTokens(user, session)
}

// AuthenticateWithPassword validates email/password and creates a session.
func (s *Service) AuthenticateWithPassword(
	ctx context.Context,
	email string,
	password string,
) (*valueobjects.AuthTokens, error) {
	// Find user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Get local account for password verification
	// Use email as the provider account ID for local auth
	account, err := s.accountRepo.GetAccountByProvider(ctx, "local", email)
	if err != nil {
		// TODO: For now, just skip password verification if no account exists
		// In production, you'd store passwords properly
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password (assuming password hash is stored in account's AccessToken field for local provider)
	// TODO: Add proper password field to Account entity or create separate local auth table
	if account.AccessToken == nil || !VerifyPassword(password, *account.AccessToken) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Create session
	session, err := s.createSession(ctx, user.ID, map[string]any{
		"auth_method": "password",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate auth tokens
	return s.generateTokens(user, session)
}

// RefreshToken validates a refresh token and issues new tokens.
func (s *Service) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*valueobjects.AuthTokens, error) {
	// Validate refresh token
	claims, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get session
	session, err := s.sessionRepo.Get(ctx, claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Get user
	user, err := s.userRepo.Get(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new auth tokens
	return s.generateTokens(user, session)
}

// GetSessionByToken retrieves a session from an access token.
func (s *Service) GetSessionByToken(
	ctx context.Context,
	accessToken string,
) (*entities.Session, error) {
	// Validate the access token
	claims, err := s.tokenManager.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	// Get the session
	session, err := s.sessionRepo.Get(ctx, claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	return session, nil
}

// ValidateAccessToken validates an access token and returns claims.
func (s *Service) ValidateAccessToken(token string) (*TokenClaims, error) {
	return s.tokenManager.ValidateAccessToken(token)
}

// Helper methods

func (s *Service) findOrCreateUser(
	ctx context.Context,
	userInfo *oauth.UserInfo,
) (*entities.User, error) {
	// Try to find existing user by email
	user, err := s.userRepo.GetUserByEmail(ctx, userInfo.Email)
	if err == nil {
		return user, nil
	}

	// Create new user
	user = &entities.User{
		ID:            uuid.New(),
		Email:         userInfo.Email,
		EmailVerified: userInfo.EmailVerified,
		Name:          userInfo.Name,
		Image:         &userInfo.Picture,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.userRepo.Create(ctx, user)
}

func (s *Service) createSession(
	ctx context.Context,
	userID uuid.UUID,
	metadata map[string]any,
) (*entities.Session, error) {
	// Extract auth method and provider from metadata
	authMethod := "local"
	authProvider := "local"

	if method, ok := metadata["auth_method"].(string); ok {
		authMethod = method
	}
	if provider, ok := metadata["provider"].(string); ok {
		authProvider = provider
		authMethod = "oauth_" + provider
	}

	// Generate a secure session token
	token := uuid.New().String()

	ipAddr := bindHost
	userAgent := "unknown"

	// Convert authProvider string to SessionAuthProvider enum
	provider := entities.SessionAuthProvider(authProvider)

	session := &entities.Session{
		ID:           uuid.New(),
		UserID:       userID,
		Token:        token,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
		AuthMethod:   &authMethod,
		AuthProvider: &provider,
		IPAddress:    &ipAddr,
		UserAgent:    &userAgent,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdSession, err := s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return createdSession, nil
}

func (s *Service) generateTokens(
	user *entities.User,
	session *entities.Session,
) (*valueobjects.AuthTokens, error) {
	// Generate access token (short-lived)
	accessToken, err := s.tokenManager.CreateAccessToken(user.ID, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate refresh token (long-lived)
	refreshToken, err := s.tokenManager.CreateRefreshToken(user.ID, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &valueobjects.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes
		SessionID:    session.ID.String(),
	}, nil
}

// Register creates a new user account and session
func (s *Service) Register(
	ctx context.Context,
	email, password, name string,
) (*valueobjects.AuthTokens, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user
	user, err := entities.NewUser(email, false, nil, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Store the user
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create a local account to store the password
	account := &entities.Account{
		ID:                uuid.New(),
		UserID:            createdUser.ID,
		Provider:          entities.AccountProviderLocal,
		AccountIdentifier: email,           // Use email as account ID for local auth
		AccessToken:       &hashedPassword, // Store password hash in AccessToken temporarily
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if _, err := s.accountRepo.Create(ctx, account); err != nil {
		// If account creation fails, delete the user
		_ = s.userRepo.Delete(ctx, createdUser.ID)
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Create a session and generate tokens
	session, err := s.createSession(ctx, createdUser.ID, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return s.generateTokens(createdUser, session)
}

// DeleteSessionByID deletes a specific session
func (s *Service) DeleteSessionByID(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

// DeleteAllUserSessions deletes all sessions for a user
func (s *Service) DeleteAllUserSessions(ctx context.Context, sessionID uuid.UUID) error {
	// Get the session to find the user ID
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}

	// Delete all sessions for the user
	// Note: We need to list all sessions and filter by user ID
	// since there's no FindByUserID method
	sessions, _, err := s.sessionRepo.List(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	for _, sess := range sessions {
		if sess.UserID == session.UserID {
			if err := s.sessionRepo.Delete(ctx, sess.ID); err != nil {
				return fmt.Errorf("failed to delete session: %w", err)
			}
		}
	}

	return nil
}

// RequestPasswordReset initiates a password reset flow
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	// Find user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists
		return nil
	}

	// Generate reset token
	token := uuid.New().String()

	// Store token in cache with expiration
	key := fmt.Sprintf("password_reset:%s", token)
	userIDStr := user.ID.String()
	if err := s.cache.Set(ctx, key, &userIDStr, 15*time.Minute); err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// TODO: Send email with reset link

	return nil
}

// ConfirmPasswordReset completes the password reset flow
func (s *Service) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	// Get user ID from token
	key := fmt.Sprintf("password_reset:%s", token)
	userIDStrPtr, err := s.cache.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	if userIDStrPtr == nil {
		return fmt.Errorf("invalid token data")
	}

	userID, err := uuid.Parse(*userIDStrPtr)
	if err != nil {
		return fmt.Errorf("invalid user ID in token")
	}

	// Hash the new password
	_, err = HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// TODO: Update password in separate password storage
	// For now, we'll just update the user's updated_at timestamp
	if _, err := s.userRepo.Update(ctx, userID, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Delete the reset token
	_ = s.cache.Delete(ctx, key)

	return nil
}

// RequestEmailVerification sends an email verification link
func (s *Service) RequestEmailVerification(ctx context.Context, sessionID uuid.UUID) error {
	// Get session to find user
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found")
	}

	user, err := s.userRepo.Get(ctx, session.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Generate verification token
	token := uuid.New().String()

	// Store token in cache
	key := fmt.Sprintf("email_verification:%s", token)
	userIDStr := user.ID.String()
	if err := s.cache.Set(ctx, key, &userIDStr, 24*time.Hour); err != nil {
		return fmt.Errorf("failed to store verification token: %w", err)
	}

	// TODO: Send email with verification link

	return nil
}

// ConfirmEmailVerification verifies the user's email
func (s *Service) ConfirmEmailVerification(ctx context.Context, token string) error {
	// Get user ID from token
	key := fmt.Sprintf("email_verification:%s", token)
	userIDStrPtr, err := s.cache.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	if userIDStrPtr == nil {
		return fmt.Errorf("invalid token data")
	}

	userID, err := uuid.Parse(*userIDStrPtr)
	if err != nil {
		return fmt.Errorf("invalid user ID in token")
	}

	// Update user as verified
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.EmailVerified = true
	if _, err := s.userRepo.Update(ctx, userID, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Delete the verification token
	_ = s.cache.Delete(ctx, key)

	return nil
}

// RequestEmailChange initiates an email change flow
func (s *Service) RequestEmailChange(
	ctx context.Context,
	sessionID uuid.UUID,
	newEmail string,
) error {
	// Get session to find user
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found")
	}

	// Check if new email is already in use
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, newEmail)
	if existingUser != nil {
		return fmt.Errorf("email already in use")
	}

	// Generate change token
	token := uuid.New().String()

	// Store token with new email
	key := fmt.Sprintf("email_change:%s", token)
	data := fmt.Sprintf("%s:%s", session.UserID.String(), newEmail)
	if err := s.cache.Set(ctx, key, &data, 15*time.Minute); err != nil {
		return fmt.Errorf("failed to store change token: %w", err)
	}

	// TODO: Send email to both old and new addresses

	return nil
}

// ConfirmEmailChange completes the email change flow
func (s *Service) ConfirmEmailChange(
	ctx context.Context,
	token, newEmail string,
	userID uuid.UUID,
) error {
	// Get data from token
	key := fmt.Sprintf("email_change:%s", token)
	dataPtr, err := s.cache.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("invalid or expired change token")
	}

	if dataPtr == nil {
		return fmt.Errorf("invalid token data")
	}

	// Parse stored data
	parts := strings.Split(*dataPtr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid token data")
	}

	storedUserID, err := uuid.Parse(parts[0])
	if err != nil || storedUserID != userID {
		return fmt.Errorf("token user mismatch")
	}

	if parts[1] != newEmail {
		return fmt.Errorf("email mismatch")
	}

	// Update user email
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.Email = newEmail
	if _, err := s.userRepo.Update(ctx, userID, user); err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	// Delete the change token
	_ = s.cache.Delete(ctx, key)

	return nil
}

// LinkAccount links an OAuth provider to an existing account
func (s *Service) LinkAccount(
	ctx context.Context,
	sessionID uuid.UUID,
	provider string,
	_ *string,
) (*entities.Account, error) {
	// Get session to find user
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	// TODO: Implement OAuth provider linking
	// This would typically involve:
	// 1. Redirecting to OAuth provider
	// 2. Handling callback
	// 3. Storing provider ID with user account

	// For now, return error - when implemented, should return the created/linked account
	_ = session
	return nil, fmt.Errorf("OAuth linking not yet implemented for provider: %s", provider)
}

// DeleteAccount deletes a user account and all sessions
func (s *Service) DeleteAccount(
	ctx context.Context,
	sessionID uuid.UUID,
) (*entities.Account, error) {
	// Get session to find user
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	// Get the account before deletion
	// TODO: Implement method to get account by user ID
	// For now, we'll need to query accounts by user ID
	// account, err := s.accountRepo.GetByUserID(ctx, session.UserID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get account: %w", err)
	// }

	// Delete all user sessions first
	if err := s.DeleteAllUserSessions(ctx, sessionID); err != nil {
		return nil, fmt.Errorf("failed to delete sessions: %w", err)
	}

	// Delete the user account
	if err := s.userRepo.Delete(ctx, session.UserID); err != nil {
		return nil, fmt.Errorf("failed to delete account: %w", err)
	}

	// Return a placeholder account for now
	// TODO: Return actual account once we can retrieve it before deletion
	account := &entities.Account{
		ID:     uuid.New(),
		UserID: session.UserID,
	}

	return account, nil
}

// UpdateAccount updates user account information
func (s *Service) UpdateAccount(
	ctx context.Context,
	sessionID uuid.UUID,
	updates map[string]any,
) (*entities.User, error) {
	// Get session to find user
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	// Get current user
	user, err := s.userRepo.Get(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Apply updates
	for key, value := range updates {
		switch key {
		case "name":
			if name, ok := value.(string); ok {
				user.Name = name
			}
		case "email":
			if email, ok := value.(string); ok {
				user.Email = email
			}
		case "image":
			if image, ok := value.(string); ok {
				user.Image = &image
			}
		case "verified":
			if verified, ok := value.(bool); ok {
				user.EmailVerified = verified
			}
		}
	}

	// Save updates
	updatedUser, err := s.userRepo.Update(ctx, session.UserID, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return updatedUser, nil
}
