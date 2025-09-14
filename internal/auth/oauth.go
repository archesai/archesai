package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
)

// OAuthService handles OAuth2 authentication flows
type OAuthService struct {
	providers      map[string]OAuthProvider
	service        *Service
	sessionManager *SessionManager
	repo           Repository
	usersRepo      users.Repository
	logger         *slog.Logger
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(
	service *Service,
	sessionManager *SessionManager,
	repo Repository,
	usersRepo users.Repository,
	logger *slog.Logger,
) *OAuthService {
	return &OAuthService{
		providers:      make(map[string]OAuthProvider),
		service:        service,
		sessionManager: sessionManager,
		repo:           repo,
		usersRepo:      usersRepo,
		logger:         logger,
	}
}

// RegisterProvider registers an OAuth provider
func (s *OAuthService) RegisterProvider(provider OAuthProvider) {
	s.providers[provider.GetProviderID()] = provider
}

// GetAuthURL generates an authorization URL for the specified provider
func (s *OAuthService) GetAuthURL(ctx context.Context, providerID string, redirectURI string) (string, string, error) {
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
	if err := s.sessionManager.StoreOAuthState(ctx, state, providerID, redirectURI, 5*time.Minute); err != nil {
		return "", "", fmt.Errorf("failed to store state: %w", err)
	}

	authURL := provider.GetAuthURL(state, redirectURI)
	return authURL, state, nil
}

// HandleCallback processes the OAuth callback
func (s *OAuthService) HandleCallback(ctx context.Context, providerID, code, state, storedState string) (*User, *TokenResponse, error) {
	// Validate state to prevent CSRF attacks
	if state != storedState {
		return nil, nil, fmt.Errorf("invalid state parameter")
	}

	provider, ok := s.providers[providerID]
	if !ok {
		return nil, nil, fmt.Errorf("provider %s not found", providerID)
	}

	// Get redirect URI from stored state
	redirectURI, err := s.sessionManager.GetOAuthRedirectURI(ctx, state)
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
	account, err := s.repo.GetAccountByProviderAndProviderID(ctx, providerID, userInfo.ProviderAccountID)
	var user *User

	if err != nil {
		// Account doesn't exist, create new user and account
		user, err = s.createOAuthUser(ctx, providerID, userInfo, tokens)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Account exists, get the user from users repository
		usersUser, err := s.usersRepo.GetUser(ctx, account.UserId)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get user: %w", err)
		}

		// Convert to auth.User
		user = &User{
			Id:            usersUser.Id,
			Email:         usersUser.Email,
			Name:          usersUser.Name,
			Image:         usersUser.Image,
			EmailVerified: usersUser.EmailVerified,
			CreatedAt:     usersUser.CreatedAt,
			UpdatedAt:     usersUser.UpdatedAt,
		}

		// Update account with new tokens if refresh token changed
		if tokens.RefreshToken != "" && tokens.RefreshToken != account.RefreshToken {
			account.RefreshToken = tokens.RefreshToken
			account.UpdatedAt = time.Now()
			if _, err := s.repo.UpdateAccount(ctx, account.Id, account); err != nil {
				s.logger.Warn("failed to update refresh token", "error", err)
			}
		}
	}

	// Convert User to users.User for token generation
	usersUser := &users.User{
		Id:            user.Id,
		Email:         user.Email,
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	// Generate JWT tokens using the service's internal method
	providerStr := providerID
	jwtTokens, err := s.service.generateTokensWithContext(
		usersUser,
		uuid.Nil,     // activeOrgID
		"",           // orgName
		"",           // orgRole
		nil,          // organizations
		nil,          // permissions
		&providerStr, // provider
		s.service.config,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Clean up OAuth state
	_ = s.sessionManager.DeleteOAuthState(ctx, state)

	return user, jwtTokens, nil
}

// createOAuthUser creates a new user from OAuth provider info
func (s *OAuthService) createOAuthUser(ctx context.Context, providerID string, userInfo *OAuthUserInfo, tokens *OAuthTokens) (*User, error) {
	now := time.Now()

	// Create user entity
	userEntity := &users.User{
		Id:            uuid.New(),
		Email:         users.Email(userInfo.Email),
		Name:          userInfo.Name,
		Image:         userInfo.Picture,
		EmailVerified: userInfo.EmailVerified,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save user to database
	createdUser, err := s.usersRepo.CreateUser(ctx, userEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create OAuth account
	account := &Account{
		Id:           uuid.New(),
		UserId:       createdUser.Id,
		ProviderId:   AccountProviderId(providerID),
		AccountId:    userInfo.ProviderAccountID,
		RefreshToken: tokens.RefreshToken,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err = s.repo.CreateAccount(ctx, account); err != nil {
		// Rollback user creation
		_ = s.usersRepo.DeleteUser(ctx, createdUser.Id)
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Convert to auth.User
	return &User{
		Id:            createdUser.Id,
		Email:         createdUser.Email,
		Name:          createdUser.Name,
		Image:         createdUser.Image,
		EmailVerified: createdUser.EmailVerified,
		CreatedAt:     createdUser.CreatedAt,
		UpdatedAt:     createdUser.UpdatedAt,
	}, nil
}

// RefreshOAuthToken refreshes an OAuth access token
func (s *OAuthService) RefreshOAuthToken(ctx context.Context, userID uuid.UUID, providerID string) (*OAuthTokens, error) {
	// Get the user's OAuth account
	accounts, _, err := s.repo.ListAccounts(ctx, ListAccountsParams{
		UserID:     &userID,
		ProviderID: &providerID,
		Limit:      1,
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

	if _, err := s.repo.UpdateAccount(ctx, account.Id, account); err != nil {
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
func (sm *SessionManager) StoreOAuthState(_ context.Context, state, provider, redirectURI string, ttl time.Duration) error {
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
func (sm *SessionManager) GetOAuthRedirectURI(_ context.Context, state string) (string, error) {
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
func (sm *SessionManager) DeleteOAuthState(_ context.Context, state string) error {
	oauthStateStore.Lock()
	defer oauthStateStore.Unlock()

	delete(oauthStateStore.states, state)
	return nil
}
