package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/sessions"
	"github.com/archesai/archesai/internal/users"
)

// Service manages OAuth providers and authentication flow
type Service struct {
	providers       map[string]Provider
	logger          *slog.Logger
	sessionsService *sessions.Service
	usersService    *users.Service
}

// NewService creates a new OAuth service with configured providers
func NewService(
	cfg *config.Config,
	sessionsService *sessions.Service,
	usersService *users.Service,
	logger *slog.Logger,
) *Service {
	s := &Service{
		providers:       make(map[string]Provider),
		logger:          logger,
		sessionsService: sessionsService,
		usersService:    usersService,
	}

	// Initialize providers based on configuration
	if cfg.Auth.Google.Enabled {
		s.providers["google"] = NewGoogleProvider(
			cfg.Auth.Google.ClientID,
			cfg.Auth.Google.ClientSecret,
			cfg.Auth.Google.RedirectURL,
			logger,
		)
		logger.Info("Google OAuth provider enabled")
	}

	if cfg.Auth.Github.Enabled {
		s.providers["github"] = NewGitHubProvider(
			cfg.Auth.Github.ClientID,
			cfg.Auth.Github.ClientSecret,
			cfg.Auth.Github.RedirectURL,
			logger,
		)
		logger.Info("GitHub OAuth provider enabled")
	}

	if cfg.Auth.Microsoft.Enabled {
		s.providers["microsoft"] = NewMicrosoftProvider(
			cfg.Auth.Microsoft.ClientID,
			cfg.Auth.Microsoft.ClientSecret,
			cfg.Auth.Microsoft.RedirectURL,
			logger,
		)
		logger.Info("Microsoft OAuth provider enabled")
	}

	return s
}

// GetAuthorizationURL generates the authorization URL for a provider
func (s *Service) GetAuthorizationURL(
	_ context.Context,
	providerName string,
	redirectURI string,
) (string, error) {
	provider, exists := s.providers[providerName]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", providerName)
	}

	// Generate a secure random state parameter
	state, err := s.generateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	// TODO: Store state in cache/database for verification against CSRF attacks
	// The state should be stored with a TTL and verified in HandleCallback

	authURL := provider.GetAuthURL(state, redirectURI)
	return authURL, nil
}

// HandleCallback processes the OAuth callback
func (s *Service) HandleCallback(
	ctx context.Context,
	providerName string,
	code string,
	state string,
	redirectURI string,
) (*sessions.TokenPair, error) {
	provider, exists := s.providers[providerName]
	if !exists {
		return nil, fmt.Errorf("provider %s not configured", providerName)
	}

	// TODO: Verify state parameter against stored value to prevent CSRF attacks
	// Should check if state exists in cache/database and hasn't expired
	_ = state // Acknowledge the parameter for now
	// For now, we'll skip state verification

	// Exchange authorization code for tokens
	tokens, err := provider.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		s.logger.Error("Failed to exchange code for tokens",
			"provider", providerName,
			"error", err)
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user information from provider
	userInfo, err := provider.GetUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		s.logger.Error("Failed to get user info",
			"provider", providerName,
			"error", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Convert to users.OAuthUserInfo to avoid circular dependency
	userOAuthInfo := &users.OAuthUserInfo{
		ProviderAccountID: userInfo.ProviderAccountID,
		Email:             userInfo.Email,
		EmailVerified:     userInfo.EmailVerified,
		Name:              userInfo.Name,
		Picture:           userInfo.Picture,
		Locale:            userInfo.Locale,
	}

	// Find or create user
	user, err := s.usersService.FindOrCreateFromOAuth(ctx, providerName, userOAuthInfo)
	if err != nil {
		s.logger.Error("Failed to find or create user from OAuth",
			"provider", providerName,
			"email", userInfo.Email,
			"error", err)
		return nil, fmt.Errorf("failed to process user: %w", err)
	}

	// Create session and generate tokens
	// For now, we'll create a simple session ID
	sessionID := uuid.New()

	accessToken, err := s.sessionsService.GenerateAccessToken(user.ID, sessionID)
	if err != nil {
		s.logger.Error("Failed to generate access token",
			"userId", user.ID,
			"error", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.sessionsService.GenerateRefreshToken(user.ID, sessionID)
	if err != nil {
		s.logger.Error("Failed to generate refresh token",
			"userId", user.ID,
			"error", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenPair := &sessions.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	s.logger.Info("OAuth login successful",
		"provider", providerName,
		"userId", user.ID,
		"email", user.Email)

	return tokenPair, nil
}

// GetProvider returns a provider by name
func (s *Service) GetProvider(name string) (Provider, bool) {
	provider, exists := s.providers[name]
	return provider, exists
}

// GetConfiguredProviders returns a list of configured provider names
func (s *Service) GetConfiguredProviders() []string {
	providers := make([]string, 0, len(s.providers))
	for name := range s.providers {
		providers = append(providers, name)
	}
	return providers
}

// BuildCallbackURL constructs the frontend callback URL with tokens or error
func (s *Service) BuildCallbackURL(
	provider string,
	tokenPair *sessions.TokenPair,
	err error,
) string {
	// Use relative path - the frontend handles the OAuth callback route
	callbackURL, _ := url.Parse("/auth/oauth/callback")
	q := callbackURL.Query()
	q.Set("provider", provider)

	if err != nil {
		q.Set("error", err.Error())
	} else if tokenPair != nil {
		q.Set("access_token", tokenPair.AccessToken)
		q.Set("refresh_token", tokenPair.RefreshToken)
	}

	callbackURL.RawQuery = q.Encode()
	return callbackURL.String()
}

// generateState generates a secure random state parameter
func (s *Service) generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
