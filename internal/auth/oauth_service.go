package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
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

	// Find or create user account
	user, err := s.linkOrCreateAccount(ctx, providerID, userInfo, tokens)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to link or create account: %w", err)
	}

	// Convert auth User to users.User for token generation
	usersUser := &users.User{
		Id:            user.Id,
		Email:         user.Email,
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	// Generate JWT tokens using the service's internal method
	jwtTokens, err := s.service.generateTokensWithContext(
		usersUser,
		uuid.Nil,
		"",
		"",
		"",
		AuthMethodOAuth,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Clean up OAuth state
	_ = s.sessionManager.DeleteOAuthState(ctx, state)

	return user, jwtTokens, nil
}

// linkOrCreateAccount finds an existing account or creates a new one
func (s *OAuthService) linkOrCreateAccount(
	ctx context.Context,
	providerID string,
	userInfo *OAuthUserInfo,
	tokens *OAuthTokens,
) (*User, error) {
	// Check if account already exists for this provider
	account, err := s.repo.GetAccountByProviderAndProviderID(ctx, providerID, userInfo.ProviderAccountID)
	if err == nil && account != nil {
		// Account exists, get the user
		user, err := s.usersRepo.GetUser(ctx, account.UserId)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		// Update tokens in account
		account.AccessToken = tokens.AccessToken
		account.AccessTokenExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
		if tokens.RefreshToken != "" {
			account.RefreshToken = tokens.RefreshToken
		}
		if tokens.IDToken != "" {
			account.IdToken = tokens.IDToken
		}

		_, err = s.repo.UpdateAccount(ctx, account.Id, account)
		if err != nil {
			s.logger.Error("failed to update account tokens", "error", err)
		}

		return &User{
			Id:            user.Id,
			Email:         user.Email,
			Name:          user.Name,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		}, nil
	}

	// Check if user exists with this email
	existingUser, _ := s.usersRepo.GetUserByEmail(ctx, userInfo.Email)

	var user *User
	if existingUser != nil {
		// User exists, link the OAuth account
		user = &User{
			Id:            existingUser.Id,
			Email:         existingUser.Email,
			Name:          existingUser.Name,
			EmailVerified: existingUser.EmailVerified,
			CreatedAt:     existingUser.CreatedAt,
			UpdatedAt:     existingUser.UpdatedAt,
		}

		// Create account link
		account = &Account{
			Id:                   uuid.New(),
			UserId:               user.Id,
			ProviderId:           AccountProviderId(providerID),
			AccountId:            userInfo.ProviderAccountID,
			AccessToken:          tokens.AccessToken,
			RefreshToken:         tokens.RefreshToken,
			IdToken:              tokens.IDToken,
			AccessTokenExpiresAt: time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second),
			Scope:                tokens.Scope,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
	} else {
		// Create new user
		newUserID := uuid.New()

		// OAuth users don't need a password since they authenticate through the provider
		newUser := &users.User{
			Id:            newUserID,
			Email:         users.Email(userInfo.Email),
			Name:          userInfo.Name,
			EmailVerified: userInfo.EmailVerified,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		createdUser, err := s.usersRepo.CreateUser(ctx, newUser)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		user = &User{
			Id:            createdUser.Id,
			Email:         createdUser.Email,
			Name:          createdUser.Name,
			EmailVerified: createdUser.EmailVerified,
			CreatedAt:     createdUser.CreatedAt,
			UpdatedAt:     createdUser.UpdatedAt,
		}

		// Create account for the new user
		account = &Account{
			Id:                   uuid.New(),
			UserId:               user.Id,
			ProviderId:           AccountProviderId(providerID),
			AccountId:            userInfo.ProviderAccountID,
			AccessToken:          tokens.AccessToken,
			RefreshToken:         tokens.RefreshToken,
			IdToken:              tokens.IDToken,
			AccessTokenExpiresAt: time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second),
			Scope:                tokens.Scope,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
	}

	// Create the account
	_, err = s.repo.CreateAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return user, nil
}

// RefreshOAuthToken refreshes an OAuth access token
func (s *OAuthService) RefreshOAuthToken(ctx context.Context, userID uuid.UUID, providerID string) (*OAuthTokens, error) {
	// Get user's account for this provider
	accounts, err := s.repo.ListUserAccounts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user accounts: %w", err)
	}

	var account *Account
	for _, acc := range accounts {
		if string(acc.ProviderId) == providerID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, fmt.Errorf("no account found for provider %s", providerID)
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

	// Update account with new tokens
	account.AccessToken = tokens.AccessToken
	account.AccessTokenExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	if tokens.RefreshToken != "" {
		account.RefreshToken = tokens.RefreshToken
	}

	_, err = s.repo.UpdateAccount(ctx, account.Id, account)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return tokens, nil
}

// UnlinkAccount removes an OAuth account link
func (s *OAuthService) UnlinkAccount(ctx context.Context, userID uuid.UUID, providerID string) error {
	// Get user's accounts
	accounts, err := s.repo.ListUserAccounts(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to list user accounts: %w", err)
	}

	// Ensure user has at least one other login method
	if len(accounts) <= 1 {
		return fmt.Errorf("cannot unlink last authentication method")
	}

	// Find and delete the account
	for _, account := range accounts {
		if string(account.ProviderId) == providerID {
			return s.repo.DeleteAccount(ctx, account.Id)
		}
	}

	return fmt.Errorf("account not found for provider %s", providerID)
}

// generateState generates a secure random state token
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// OAuthSessionManager adds OAuth-specific methods to SessionManager
type OAuthSessionManager interface {
	StoreOAuthState(ctx context.Context, state, provider, redirectURI string, ttl time.Duration) error
	GetOAuthRedirectURI(ctx context.Context, state string) (string, error)
	DeleteOAuthState(ctx context.Context, state string) error
}

// OAuth state methods are implemented in oauth_state.go
