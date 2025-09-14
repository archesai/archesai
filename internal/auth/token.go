package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT token claims (legacy - use EnhancedClaims instead)
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// TokenResponseWithExpiry extends the generated TokenResponse with ExpiresAt
type TokenResponseWithExpiry struct {
	TokenResponse
	ExpiresAt time.Time `json:"expires_at"`
}

// ValidateToken validates and parses a JWT token using enhanced claims
func (s *Service) ValidateToken(tokenString string) (*EnhancedClaims, error) {
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
		return nil, ErrInvalidToken
	}

	// Additional validation for token type
	if claims.TokenType != AccessTokenType {
		return nil, fmt.Errorf("invalid token type: expected access token")
	}

	return claims, nil
}

// ValidateLegacyToken validates and parses a JWT token using legacy claims
// Deprecated: Use ValidateToken instead
func (s *Service) ValidateLegacyToken(tokenString string) (*Claims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAPIKey validates an API key
func (s *Service) ValidateAPIKey(ctx context.Context, apiKey string) (*APIKey, error) {
	if s.apiKeyService == nil {
		return nil, fmt.Errorf("API key service not configured")
	}

	// Delegate to API key service
	return s.apiKeyService.ValidateAPIKey(ctx, apiKey)
}

// RefreshToken refreshes an access token using a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Parse the refresh token with enhanced claims
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
		return nil, ErrInvalidToken
	}

	// Verify token type
	if claims.TokenType != RefreshTokenType {
		return nil, fmt.Errorf("invalid token type: expected refresh token")
	}

	// Get the user
	user, err := s.usersRepo.GetUser(ctx, claims.UserID)
	if err != nil {
		s.logger.Error("failed to get user for refresh", "user_id", claims.UserID, "error", err)
		return nil, ErrUserNotFound
	}

	// Verify session is still valid if sessionID is present
	if claims.SessionID != "" {
		sessionID, err := uuid.Parse(claims.SessionID)
		if err == nil {
			session, err := s.repo.GetSession(ctx, sessionID)
			if err != nil || session == nil {
				s.logger.Warn("session not found for refresh", "session_id", claims.SessionID)
				return nil, ErrSessionNotFound
			}

			// Check if session is expired
			if session.ExpiresAt != "" {
				expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
				if err == nil && time.Now().After(expiresAt) {
					s.logger.Warn("session expired for refresh", "session_id", claims.SessionID)
					return nil, ErrSessionExpired
				}
			}
		}
	}

	// Generate new tokens
	newTokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	s.logger.Info("tokens refreshed successfully", "user_id", claims.UserID)
	return newTokens, nil
}

// generateTokens generates access and refresh tokens for a user
func (s *Service) generateTokens(user *users.User) (*TokenResponse, error) {
	return s.generateTokensWithConfig(user, s.config)
}

// generateTokensWithConfig generates tokens with a specific configuration
func (s *Service) generateTokensWithConfig(user *users.User, config Config) (*TokenResponse, error) {
	// TODO: Fetch user's organizations and permissions
	// For now, use empty slices
	var organizations []OrganizationClaim
	var permissions []string

	return s.generateTokensWithContext(
		user,
		uuid.Nil, // No active organization for now
		"",       // No organization name
		"",       // No organization role
		organizations,
		permissions,
		nil, // No provider for local auth
		config,
	)
}

// generateTokensWithContext generates tokens with full context
func (s *Service) generateTokensWithContext(
	user *users.User,
	activeOrgID uuid.UUID,
	orgName string,
	orgRole string,
	organizations []OrganizationClaim,
	permissions []string,
	provider *string,
	config Config,
) (*TokenResponse, error) {
	now := time.Now()

	// Create enhanced claims for access token
	accessClaims := &EnhancedClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.Id.String(),
		},
		UserID:           user.Id,
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
			ExpiresAt: jwt.NewNumericDate(now.Add(config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.Id.String(),
		},
		UserID:     user.Id,
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
		ExpiresIn:    int64(config.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}
