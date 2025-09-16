// Package oauth provides OAuth2 authentication services.
package oauth

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package oauth --include-tags OAuth ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package oauth --include-tags OAuth ../../api/openapi.bundled.yaml

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/sessions"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
)

// Service handles OAuth2 authentication flows
type Service struct {
	providers      map[string]Provider
	sessionManager *sessions.SessionManager
	repo           accounts.Repository
	usersRepo      users.Repository
	logger         *slog.Logger
}

// NewService creates a new OAuth service
func NewService(
	sessionManager *sessions.SessionManager,
	repo accounts.Repository,
	usersRepo users.Repository,
	logger *slog.Logger,
) *Service {
	return &Service{
		providers:      make(map[string]Provider),
		sessionManager: sessionManager,
		repo:           repo,
		usersRepo:      usersRepo,
		logger:         logger,
	}
}

// RegisterProvider registers an OAuth provider
func (s *Service) RegisterProvider(provider Provider) {
	s.providers[provider.GetProviderID()] = provider
}

// GetAuthURL generates an authorization URL for the specified provider
func (s *Service) GetAuthURL(ctx context.Context, providerID string, redirectURI string) (string, string, error) {
	provider, ok := s.providers[providerID]
	if !ok {
		return "", "", fmt.Errorf("provider %s not found", providerID)
	}

	// Generate secure state token
	state, err := generateState()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Store state in session with expiry (5 minutes)
	if err := s.StoreOAuthState(ctx, state, providerID, redirectURI, 5*time.Minute); err != nil {
		return "", "", fmt.Errorf("failed to store state: %w", err)
	}

	authURL := provider.GetAuthURL(state, redirectURI)
	return authURL, state, nil
}

// HandleCallback processes the OAuth callback
func (s *Service) HandleCallback(ctx context.Context, providerID, code, state, storedState string) (*users.User, *sessions.TokenResponse, error) {
	// Validate state to prevent CSRF attacks
	if state != storedState {
		return nil, nil, fmt.Errorf("invalid state parameter")
	}

	provider, ok := s.providers[providerID]
	if !ok {
		return nil, nil, fmt.Errorf("provider %s not found", providerID)
	}

	// Get redirect URI from stored state
	redirectURI, err := s.GetOAuthRedirectURI(ctx, state)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get redirect URI: %w", err)
	}

	// Exchange code for tokens
	tokens, err := provider.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from provider
	userInfo, err := provider.GetUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user already exists
	account, err := s.repo.GetByProviderID(ctx, providerID, userInfo.ProviderAccountID)
	var user *users.User

	if err != nil {
		// Account doesn't exist, create new user and account
		user, err = s.createOAuthUser(ctx, providerID, userInfo, tokens)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Account exists, get the user from users repository
		usersUser, err := s.usersRepo.Get(ctx, account.UserID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get user: %w", err)
		}

		// Use the existing users.User
		user = usersUser

		// Update account with new tokens if refresh token changed
		if tokens.RefreshToken != "" && tokens.RefreshToken != account.RefreshToken {
			account.RefreshToken = tokens.RefreshToken
			account.UpdatedAt = time.Now()
			if _, err := s.repo.Update(ctx, account.ID, account); err != nil {
				s.logger.Warn("failed to update refresh token", "error", err)
			}
		}
	}

	// Create a session for the authenticated user
	// Note: This assumes the user has a default organization
	// In a real implementation, you'd need to determine the appropriate organization
	session, err := s.sessionManager.Create(ctx, user.ID, uuid.New(), "", "")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create token response using session token
	tokenResponse := &sessions.TokenResponse{
		AccessToken:  session.Token,
		RefreshToken: session.Token, // For simplicity, using same token
		TokenType:    "bearer",
		ExpiresIn:    int64(24 * 60 * 60), // 24 hours
	}

	// Clean up OAuth state
	_ = s.DeleteOAuthState(ctx, state)

	return user, tokenResponse, nil
}

// createOAuthUser creates a new user from OAuth provider info
func (s *Service) createOAuthUser(ctx context.Context, providerID string, userInfo *UserInfo, tokens *Tokens) (*users.User, error) {
	now := time.Now()

	// Create user entity
	userEntity := &users.User{
		ID:            uuid.New(),
		Email:         users.Email(userInfo.Email),
		Name:          userInfo.Name,
		Image:         userInfo.Picture,
		EmailVerified: userInfo.EmailVerified,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save user to database
	createdUser, err := s.usersRepo.Create(ctx, userEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create OAuth account
	account := &accounts.Account{
		ID:           uuid.New(),
		UserID:       createdUser.ID,
		ProviderID:   accounts.AccountProviderID(providerID),
		AccountID:    userInfo.ProviderAccountID,
		RefreshToken: tokens.RefreshToken,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err := s.repo.Create(ctx, account); err != nil {
		// Rollback user creation
		_ = s.usersRepo.Delete(ctx, createdUser.ID)
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Return the created users.User
	return createdUser, nil
}

// RefreshOAuthToken refreshes an OAuth access token
func (s *Service) RefreshOAuthToken(ctx context.Context, _ uuid.UUID, providerID string) (*Tokens, error) {
	// Get the user's OAuth account
	accounts, _, err := s.repo.List(ctx, accounts.ListAccountsParams{
		Page: accounts.PageQuery{
			Number: 1,
			Size:   1,
		},
	})
	if err != nil || len(accounts) == 0 {
		return nil, fmt.Errorf("OAuth account not found")
	}

	account := accounts[0]
	if account.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	provider, ok := s.providers[providerID]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", providerID)
	}

	// Refresh the token
	tokens, err := provider.RefreshToken(ctx, account.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Update the account with new tokens
	account.RefreshToken = tokens.RefreshToken
	account.UpdatedAt = time.Now()

	if _, err := s.repo.Update(ctx, account.ID, account); err != nil {
		s.logger.Warn("failed to update account tokens", "error", err)
	}

	return tokens, nil
}

// generateState generates a secure random state parameter
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// OAuth state storage using a simple in-memory map
// In production, this should use Redis or another distributed cache
var oauthStateStore = struct {
	sync.RWMutex
	states map[string]oauthState
}{
	states: make(map[string]oauthState),
}

type oauthState struct {
	Provider    string
	RedirectURI string
	ExpiresAt   time.Time
}

// StoreOAuthState stores OAuth state data temporarily
func (s *Service) StoreOAuthState(_ context.Context, state, provider, redirectURI string, ttl time.Duration) error {
	oauthStateStore.Lock()
	defer oauthStateStore.Unlock()

	// Clean up expired states
	now := time.Now()
	for k, v := range oauthStateStore.states {
		if v.ExpiresAt.Before(now) {
			delete(oauthStateStore.states, k)
		}
	}

	oauthStateStore.states[state] = oauthState{
		Provider:    provider,
		RedirectURI: redirectURI,
		ExpiresAt:   now.Add(ttl),
	}

	return nil
}

// GetOAuthRedirectURI retrieves the redirect URI for a state
func (s *Service) GetOAuthRedirectURI(_ context.Context, state string) (string, error) {
	oauthStateStore.RLock()
	defer oauthStateStore.RUnlock()

	if s, ok := oauthStateStore.states[state]; ok {
		if s.ExpiresAt.After(time.Now()) {
			return s.RedirectURI, nil
		}
	}

	return "", fmt.Errorf("state not found or expired")
}

// DeleteOAuthState removes OAuth state data
func (s *Service) DeleteOAuthState(_ context.Context, state string) error {
	oauthStateStore.Lock()
	defer oauthStateStore.Unlock()

	delete(oauthStateStore.states, state)
	return nil
}
