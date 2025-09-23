package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/users"
)

// Provider represents an OAuth provider
type Provider string

// OAuth providers
const (
	ProviderGoogle    Provider = "google"
	ProviderGitHub    Provider = "github"
	ProviderMicrosoft Provider = "microsoft"
	ProviderApple     Provider = "apple"
	ProviderLocal     Provider = "local"
)

// Service is the unified authentication service
type Service struct {
	config          *config.Config
	logger          *slog.Logger
	usersService    *users.Service
	accountsRepo    accounts.Repository
	sessionsStore   SessionStore
	tokenManager    TokenManager
	magicLinkStore  MagicLinkStore
	apiKeyStore     APIKeyStore
	apiKeyValidator APIKeyValidator
	oauthProviders  map[Provider]OAuthProvider
	deliverers      map[DeliveryMethod]Deliverer
}

// SessionStore manages session persistence
type SessionStore interface {
	Create(ctx context.Context, userID uuid.UUID, metadata map[string]interface{}) (*Session, error)
	Get(ctx context.Context, sessionID uuid.UUID) (*Session, error)
	Delete(ctx context.Context, sessionID uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]*Session, error)
}

// TokenManager handles token generation and validation
type TokenManager interface {
	GenerateAccessToken(userID, sessionID uuid.UUID) (string, error)
	GenerateRefreshToken(userID, sessionID uuid.UUID) (string, error)
	GenerateAPIKey(userID uuid.UUID, name string, scopes []string) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
	RefreshToken(refreshToken string) (string, string, error)
}

// MagicLinkStore manages magic link tokens
type MagicLinkStore interface {
	CreateToken(
		ctx context.Context,
		identifier string,
		deliveryMethod DeliveryMethod,
		userID *uuid.UUID,
		IPAddress string,
		userAgent string,
	) (*MagicLinkToken, error)
	VerifyToken(ctx context.Context, token string) (*MagicLinkToken, error)
	VerifyOTP(ctx context.Context, identifier string, code string) (*MagicLinkToken, error)
	CleanupExpired(ctx context.Context) error
}

// OAuthProvider interface for OAuth providers
type OAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*OAuthTokens, error)
	GetUserInfo(ctx context.Context, accessToken string) (*OAuthUserInfo, error)
	RefreshToken(ctx context.Context, refreshToken string) (*OAuthTokens, error)
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID           uuid.UUID
	SessionID        uuid.UUID
	Email            string
	Name             string
	Picture          string
	Provider         string
	ProviderID       string
	OrganizationName string
	OrganizationRole string
	Roles            []string
	Permissions      []string
	Scopes           []string
	ExpiresAt        time.Time
}

// MagicLinkToken represents a magic link token
type MagicLinkToken struct {
	ID             uuid.UUID
	Token          string
	Code           string
	Identifier     string // Email or phone
	UserID         *uuid.UUID
	DeliveryMethod DeliveryMethod
	ExpiresAt      time.Time
	UsedAt         *time.Time
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time
}

// OAuthTokens represents OAuth tokens
type OAuthTokens struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	ExpiresIn    int
	Scope        string
}

// OAuthUserInfo represents user info from OAuth provider
type OAuthUserInfo struct {
	ID            string
	Email         string
	Name          string
	Picture       string
	EmailVerified bool
	Provider      Provider
}

// NewService creates a new unified auth service
func NewService(
	cfg *config.Config,
	logger *slog.Logger,
	usersService *users.Service,
	accountsRepo accounts.Repository,
	sessionsStore SessionStore,
	tokenManager TokenManager,
	magicLinkStore MagicLinkStore,
	apiKeyStore APIKeyStore,
	apiKeyValidator APIKeyValidator,
) *Service {
	s := &Service{
		config:          cfg,
		logger:          logger,
		usersService:    usersService,
		accountsRepo:    accountsRepo,
		sessionsStore:   sessionsStore,
		tokenManager:    tokenManager,
		magicLinkStore:  magicLinkStore,
		apiKeyStore:     apiKeyStore,
		apiKeyValidator: apiKeyValidator,
		oauthProviders:  make(map[Provider]OAuthProvider),
		deliverers:      make(map[DeliveryMethod]Deliverer),
	}

	// Default deliverers will be registered in app initialization

	return s
}

// RegisterProvider registers an OAuth provider
func (s *Service) RegisterProvider(provider Provider, impl OAuthProvider) {
	s.oauthProviders[provider] = impl
}

// RegisterDeliverer registers a magic link deliverer
func (s *Service) RegisterDeliverer(method DeliveryMethod, impl Deliverer) {
	s.deliverers[method] = impl
}

// GetOAuthURL returns the OAuth authentication URL for a provider
func (s *Service) GetOAuthURL(provider Provider) string {
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return ""
	}

	// Generate state for CSRF protection
	state := generateRandomState()
	return oauthProvider.GetAuthURL(state)
}

func generateRandomState() string {
	// Would implement secure random state generation
	return "random-state"
}

// AuthenticateWithPassword authenticates with username/password
func (s *Service) AuthenticateWithPassword(
	_ context.Context,
	email, password string,
) (*Session, error) {
	// TODO: Password authentication requires direct database access
	// The accounts repository doesn't expose password fields in domain model
	// Need to either:
	// 1. Add password handling to accounts repository
	// 2. Create a separate password store
	// 3. Access database directly for password operations

	// For now, return not implemented
	_ = email
	_ = password
	return nil, fmt.Errorf("password authentication not yet fully implemented")
}

// AuthenticateWithMagicLink sends a magic link
func (s *Service) AuthenticateWithMagicLink(ctx context.Context, email string) error {
	// Find or create user
	user, err := s.usersService.GetByEmail(ctx, email)
	if err != nil {
		// User doesn't exist, create them
		user, err = s.usersService.Create(ctx, &users.CreateUserRequest{
			Email:         email,
			Name:          "",    // Will be filled in later if provided
			EmailVerified: false, // Will be verified through magic link
		})
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate magic link token using user ID
	userPtr := &user.ID
	token, err := s.magicLinkStore.CreateToken(ctx, email, DeliveryConsole, userPtr, "", "")
	if err != nil {
		return fmt.Errorf("failed to create magic link: %w", err)
	}

	// Send email (would be implemented with email service)
	s.logger.Info("Magic link created", "email", email, "token", token.Token)

	return nil
}

// VerifyMagicLink verifies a magic link token
func (s *Service) VerifyMagicLink(ctx context.Context, token string) (*Session, error) {
	// Verify token
	magicLink, err := s.magicLinkStore.VerifyToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token")
	}

	// Find user
	user, err := s.usersService.GetByEmail(ctx, magicLink.Identifier)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Token is automatically marked as used during verification

	// Create session
	return s.createSession(ctx, user.ID, string(AuthMethodMagicLink), ProviderLocal, nil)
}

// AuthenticateWithOAuth handles OAuth authentication
func (s *Service) AuthenticateWithOAuth(
	ctx context.Context,
	provider Provider,
	code string,
) (*Session, error) {
	// Get OAuth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	// Exchange code for tokens
	tokens, err := oauthProvider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from provider
	userInfo, err := oauthProvider.GetUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Find or create user
	user, err := s.usersService.FindOrCreateFromOAuth(ctx, string(provider), &users.OAuthUserInfo{
		ProviderAccountID: userInfo.ID,
		Email:             userInfo.Email,
		Name:              userInfo.Name,
		Picture:           userInfo.Picture,
		EmailVerified:     userInfo.EmailVerified,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process user: %w", err)
	}

	// Create or update account record
	account := &accounts.Account{
		ID:           uuid.New(),
		UserID:       user.ID,
		AccountID:    userInfo.ID,
		ProviderID:   s.mapProviderToAccountProvider(provider),
		AccessToken:  &tokens.AccessToken,
		RefreshToken: &tokens.RefreshToken,
		IDToken:      &tokens.IDToken,
		Scope:        &tokens.Scope,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = s.accountsRepo.Create(ctx, account)
	if err != nil {
		s.logger.Error("Failed to create account record", "error", err)
		// Continue anyway - session is more important
	}

	// Create session
	metadata := map[string]interface{}{
		"oauth_tokens": tokens,
	}
	return s.createSession(ctx, user.ID, string(AuthMethodOAuth), provider, metadata)
}

// SetPassword sets or updates a user's password
func (s *Service) SetPassword(_ context.Context, userID uuid.UUID, password string) error {
	// TODO: Password management requires direct database access
	// The accounts repository doesn't expose password fields in domain model
	// Need to either:
	// 1. Add password handling to accounts repository
	// 2. Create a separate password store
	// 3. Access database directly for password operations

	// For now, return not implemented
	_ = userID
	_ = password
	return fmt.Errorf("password management not yet fully implemented")
}

// CreateAPIKey creates a new API token
func (s *Service) CreateAPIKey(
	_ context.Context,
	userID uuid.UUID,
	name string,
	scopes []string,
) (string, error) {
	// Generate API token
	token, err := s.tokenManager.GenerateAPIKey(userID, name, scopes)
	if err != nil {
		return "", fmt.Errorf("failed to generate API token: %w", err)
	}

	// Store token record (would be implemented with token repository)
	s.logger.Info("API token created", "userID", userID, "name", name)

	return token, nil
}

// ValidateToken validates any token (access, refresh, or API)
func (s *Service) ValidateToken(_ context.Context, token string) (*TokenClaims, error) {
	return s.tokenManager.ValidateToken(token)
}

// RefreshSession refreshes a session using refresh token
func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (*Session, error) {
	// Validate and refresh tokens
	newAccessToken, newRefreshToken, err := s.tokenManager.RefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Get claims from new access token
	claims, err := s.tokenManager.ValidateToken(newAccessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get existing session
	session, err := s.sessionsStore.Get(ctx, claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Update tokens
	session.AccessToken = newAccessToken
	session.RefreshToken = newRefreshToken
	session.UpdatedAt = time.Now()

	return session, nil
}

// RevokeSession revokes a session
func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionsStore.Delete(ctx, sessionID)
}

// RevokeAllSessions revokes all sessions for a user
func (s *Service) RevokeAllSessions(ctx context.Context, userID uuid.UUID) error {
	return s.sessionsStore.DeleteByUserID(ctx, userID)
}

// Helper methods

func (s *Service) createSession(
	ctx context.Context,
	userID uuid.UUID,
	method string,
	provider Provider,
	metadata map[string]interface{},
) (*Session, error) {
	// Add auth method and provider to metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["auth_method"] = method
	metadata["auth_provider"] = string(provider)

	// Get IP address and user agent from context if available
	if ctx.Value("ip_address") != nil {
		metadata["ip_address"] = ctx.Value("ip_address").(string)
	}
	if ctx.Value("user_agent") != nil {
		metadata["user_agent"] = ctx.Value("user_agent").(string)
	}

	// Store session first to get ID
	stored, err := s.sessionsStore.Create(ctx, userID, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	// Generate tokens using the stored session ID
	accessToken, err := s.tokenManager.GenerateAccessToken(userID, stored.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(userID, stored.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Build complete session response
	session := &Session{
		ID:             stored.ID,
		UserID:         userID,
		Token:          stored.Token,
		OrganizationID: stored.OrganizationID,
		AuthProvider:   string(provider),
		AuthMethod:     &method,
		IpAddress:      stored.IpAddress,
		UserAgent:      stored.UserAgent,
		ExpiresAt:      stored.ExpiresAt,
		CreatedAt:      stored.CreatedAt,
		UpdatedAt:      stored.UpdatedAt,
	}

	return session, nil
}

func (s *Service) mapProviderToAccountProvider(provider Provider) accounts.AccountProviderID {
	switch provider {
	case ProviderGoogle:
		return accounts.Google
	case ProviderGitHub:
		return accounts.Github
	case ProviderMicrosoft:
		return accounts.Microsoft
	case ProviderApple:
		return accounts.Apple
	default:
		return accounts.Local
	}
}

// APIKeyStore manages API token persistence
type APIKeyStore interface {
	CreateToken(
		ctx context.Context,
		userID, organizationID uuid.UUID,
		name string,
		scopes []string,
		expiresIn time.Duration,
	) (*APIKey, error)
	ValidateToken(ctx context.Context, key string) (*APIKey, error)
	RevokeToken(ctx context.Context, keyID uuid.UUID) error
	ListTokensByUser(ctx context.Context, userID uuid.UUID) ([]*APIKey, error)
	ParseAPIKey(authHeader string) string
}

// APIKeyValidator validates API tokens with rate limiting and scope checking
type APIKeyValidator interface {
	ValidateAPIKey(ctx context.Context, key string) (*APIKey, error)
	ValidateAPIKeyWithScopes(
		ctx context.Context,
		key string,
		requiredScopes []string,
	) (*APIKey, error)
	ValidateAPIKeyForOrganization(
		ctx context.Context,
		key string,
		organizationID uuid.UUID,
	) (*APIKey, error)
	ExtractAPIKeyFromHeaders(headers map[string]string) string
	CheckRateLimit(ctx context.Context, token *APIKey, window time.Duration) error
	ValidateScopes(tokenScopes, requiredScopes []string) error
}

// DeliveryMethod represents how magic links are delivered
type DeliveryMethod string

// Delivery methods for magic links
const (
	DeliveryEmail   DeliveryMethod = "email"
	DeliveryConsole DeliveryMethod = "console"
	DeliveryOTP     DeliveryMethod = "otp"
	DeliveryWebhook DeliveryMethod = "webhook"
)

// Deliverer delivers magic links via various methods
type Deliverer interface {
	Deliver(ctx context.Context, token *MagicLinkToken, baseURL string) error
}

// TokenPair represents an access/refresh token pair
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
}

// Error definitions
var (
	ErrSessionNotFound          = fmt.Errorf("session not found")
	ErrSessionExpired           = fmt.Errorf("session expired")
	ErrTokenAlreadyUsed         = fmt.Errorf("token already used")
	ErrTokenExpired             = fmt.Errorf("token expired")
	ErrInvalidOTP               = fmt.Errorf("invalid OTP")
	ErrRateLimitExceeded        = fmt.Errorf("rate limit exceeded")
	ErrInvalidAPIKey            = fmt.Errorf("invalid API key")
	ErrInvalidAPIKeyFormat      = fmt.Errorf("invalid API key format")
	ErrAPIKeyExpired            = fmt.Errorf("API key expired")
	ErrInsufficientScopes       = fmt.Errorf("insufficient scopes")
	ErrUnauthorizedOrganization = fmt.Errorf("unauthorized organization")
)
