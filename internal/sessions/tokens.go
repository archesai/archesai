package sessions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenConfig holds JWT token configuration
type TokenConfig struct {
	JWTSecret             string
	AccessTokenExpiry     time.Duration
	RefreshTokenExpiry    time.Duration
	SessionTokenExpiry    time.Duration
	MaxConcurrentSessions int
}

// TokenService handles JWT token operations
type TokenService struct {
	config       TokenConfig
	jwtSecret    []byte
	sessionsRepo Repository
	usersRepo    users.Repository
}

// NewTokenService creates a new JWT token service
func NewTokenService(config TokenConfig, sessionsRepo Repository, usersRepo users.Repository) *TokenService {
	// Set default values
	if config.AccessTokenExpiry == 0 {
		config.AccessTokenExpiry = 15 * time.Minute
	}
	if config.RefreshTokenExpiry == 0 {
		config.RefreshTokenExpiry = 7 * 24 * time.Hour
	}
	if config.SessionTokenExpiry == 0 {
		config.SessionTokenExpiry = 30 * 24 * time.Hour
	}

	return &TokenService{
		config:       config,
		jwtSecret:    []byte(config.JWTSecret),
		sessionsRepo: sessionsRepo,
		usersRepo:    usersRepo,
	}
}

// ValidateToken validates and parses a JWT token using enhanced claims
func (s *TokenService) ValidateToken(tokenString string) (*EnhancedClaims, error) {
	// Remove Bearer prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	// Parse and validate the token with enhanced claims
	token, err := jwt.ParseWithClaims(tokenString, &EnhancedClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*EnhancedClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Additional validation for token type
	if claims.TokenType != AccessTokenType {
		return nil, fmt.Errorf("invalid token type: expected access token")
	}

	return claims, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Parse the refresh token with refresh claims
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verify token type
	if claims.TokenType != RefreshTokenType {
		return nil, fmt.Errorf("invalid token type: expected refresh token")
	}

	// Get the user
	user, err := s.usersRepo.Get(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Verify session is still valid if sessionID is present
	if claims.SessionID != "" {
		sessionID, err := uuid.Parse(claims.SessionID)
		if err == nil {
			session, err := s.sessionsRepo.Get(ctx, sessionID)
			if err != nil || session == nil {
				return nil, fmt.Errorf("session not found")
			}

			// Check if session is expired
			// ExpiresAt is now time.Time
			if time.Now().After(session.ExpiresAt) {
				return nil, fmt.Errorf("session expired")
			}
		}
	}

	// Generate new tokens
	newTokens, err := s.GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return newTokens, nil
}

// GenerateTokens generates access and refresh tokens for a user
func (s *TokenService) GenerateTokens(user *users.User) (*TokenResponse, error) {
	return s.GenerateTokensWithContext(
		user,
		uuid.Nil, // No active organization for now
		"",       // No organization name
		"",       // No organization role
		nil,      // No organizations
		nil,      // No permissions
		nil,      // No provider for local auth
	)
}

// GenerateTokensWithContext generates tokens with full context
func (s *TokenService) GenerateTokensWithContext(
	user *users.User,
	activeOrgID uuid.UUID,
	orgName string,
	orgRole string,
	organizations []OrganizationClaim,
	permissions []string,
	provider *string,
) (*TokenResponse, error) {
	now := time.Now()

	// Create enhanced claims for access token
	accessClaims := &EnhancedClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
		UserID:           user.ID,
		Email:            string(user.Email),
		Name:             user.Name,
		OrganizationID:   activeOrgID,
		OrganizationName: orgName,
		OrganizationRole: orgRole,
		Organizations:    organizations,
		Permissions:      permissions,
		TokenType:        AccessTokenType,
		AuthMethod:       AuthMethodPassword,
		EmailVerified:    user.EmailVerified,
	}

	// Add provider info if OAuth
	if provider != nil && *provider != "" {
		accessClaims.Provider = *provider
		accessClaims.AuthMethod = AuthMethodOAuth
	}

	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Create refresh token with minimal claims
	refreshClaims := &RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
		UserID:     user.ID,
		TokenType:  RefreshTokenType,
		AuthMethod: accessClaims.AuthMethod,
	}

	// Create refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.config.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}
