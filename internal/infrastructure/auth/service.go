package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/core/repositories"
	"github.com/archesai/archesai/internal/infrastructure/auth/oauth"
	"github.com/archesai/archesai/internal/infrastructure/config"
)

// Service provides authentication functionality across all transport layers.
type Service struct {
	config         *config.Config
	sessionRepo    repositories.SessionRepository
	userRepo       repositories.UserRepository
	tokenManager   *TokenManager
	magicLink      *MagicLinkProvider
	oauthProviders map[string]OAuthProvider
}

// OAuthProvider interface for OAuth providers
type OAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth.OAuthTokens, error)
	GetUserInfo(ctx context.Context, accessToken string) (*oauth.OAuthUserInfo, error)
}

// NewService creates a new authentication service.
func NewService(
	cfg *config.Config,
	sessionRepo repositories.SessionRepository,
	userRepo repositories.UserRepository,
) *Service {
	// Build base URL from API config
	baseURL := fmt.Sprintf("http://%s:%d", cfg.API.Host, int(cfg.API.Port))
	if cfg.API.Host == "0.0.0.0" {
		baseURL = fmt.Sprintf("http://localhost:%d", int(cfg.API.Port))
	}

	s := &Service{
		config:         cfg,
		sessionRepo:    sessionRepo,
		userRepo:       userRepo,
		tokenManager:   NewTokenManager(cfg.Auth.Local.JWTSecret),
		magicLink:      NewMagicLinkProvider(cfg.Auth.Local.JWTSecret, baseURL),
		oauthProviders: make(map[string]OAuthProvider),
	}

	// Initialize OAuth providers based on config
	if cfg.Auth.Google != nil && cfg.Auth.Google.Enabled && cfg.Auth.Google.ClientId != nil {
		s.oauthProviders["google"] = oauth.NewGoogleProvider(
			*cfg.Auth.Google.ClientId,
			*cfg.Auth.Google.ClientSecret,
			*cfg.Auth.Google.RedirectUrl,
		)
	}

	if cfg.Auth.Github != nil && cfg.Auth.Github.Enabled && cfg.Auth.Github.ClientId != nil {
		s.oauthProviders["github"] = oauth.NewGitHubProvider(
			*cfg.Auth.Github.ClientId,
			*cfg.Auth.Github.ClientSecret,
			*cfg.Auth.Github.RedirectUrl,
		)
	}

	if cfg.Auth.Microsoft != nil && cfg.Auth.Microsoft.Enabled &&
		cfg.Auth.Microsoft.ClientId != nil {
		s.oauthProviders["microsoft"] = oauth.NewMicrosoftProvider(
			*cfg.Auth.Microsoft.ClientId,
			*cfg.Auth.Microsoft.ClientSecret,
			*cfg.Auth.Microsoft.RedirectUrl,
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
	state string,
) (*AuthTokens, error) {
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
	session, err := s.createSession(ctx, user.ID, map[string]interface{}{
		"provider":    provider,
		"provider_id": userInfo.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate auth tokens
	return s.generateAuthTokens(user, session)
}

// GenerateMagicLink creates a stateless magic link token.
func (s *Service) GenerateMagicLink(identifier string, redirectURL string) (string, error) {
	return s.magicLink.GenerateLink(identifier, redirectURL)
}

// VerifyMagicLink validates a magic link token and creates a session.
func (s *Service) VerifyMagicLink(ctx context.Context, token string) (*AuthTokens, error) {
	// Validate the magic link token
	claims, err := s.magicLink.ValidateLink(token)
	if err != nil {
		return nil, fmt.Errorf("invalid magic link: %w", err)
	}

	// Find or create user by identifier (email)
	user, err := s.userRepo.GetByEmail(ctx, claims.Identifier)
	if err != nil {
		// Create new user if not exists
		user, err = s.userRepo.Create(ctx, &entities.User{
			ID:            uuid.New(),
			Email:         claims.Identifier,
			EmailVerified: true,              // Magic link verifies email
			Name:          claims.Identifier, // Default to email
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Create session
	session, err := s.createSession(ctx, user.ID, map[string]interface{}{
		"auth_method": "magic_link",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate auth tokens
	return s.generateAuthTokens(user, session)
}

// RefreshToken validates a refresh token and issues new tokens.
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error) {
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
	return s.generateAuthTokens(user, session)
}

// ValidateAccessToken validates an access token and returns claims.
func (s *Service) ValidateAccessToken(token string) (*TokenClaims, error) {
	return s.tokenManager.ValidateAccessToken(token)
}

// Helper methods

func (s *Service) findOrCreateUser(
	ctx context.Context,
	userInfo *oauth.OAuthUserInfo,
) (*entities.User, error) {
	// Try to find existing user by email
	user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
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
	metadata map[string]interface{},
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

	session := &entities.Session{
		ID:           uuid.New(),
		UserID:       userID,
		Token:        token,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
		AuthMethod:   &authMethod,
		AuthProvider: &authProvider,
		IpAddress:    "0.0.0.0", // TODO: Get from request context
		UserAgent:    "unknown", // TODO: Get from request context
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.sessionRepo.Create(ctx, session)
}

func (s *Service) generateAuthTokens(
	user *entities.User,
	session *entities.Session,
) (*AuthTokens, error) {
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

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes
		SessionID:    session.ID.String(),
	}, nil
}

// AuthTokens represents the tokens returned after successful authentication.
type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	SessionID    string `json:"sessionId"`
}
