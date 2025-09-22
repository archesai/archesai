package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/auth"
)

// Manager handles JWT token generation and validation.
type Manager struct {
	jwtSecret string
}

// NewManager creates a new token manager.
func NewManager(jwtSecret string) *Manager {
	return &Manager{
		jwtSecret: jwtSecret,
	}
}

// GenerateTokenPair generates both access and refresh tokens.
func (m *Manager) GenerateTokenPair(
	userID uuid.UUID,
	sessionID uuid.UUID,
	orgID uuid.UUID,
	claims *auth.TokenClaims,
) (*auth.TokenPair, error) {
	accessToken, err := m.GenerateAccessTokenWithClaims(userID, sessionID, orgID, claims)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := m.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &auth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// GenerateAccessToken generates a JWT access token (interface method).
func (m *Manager) GenerateAccessToken(userID, sessionID uuid.UUID) (string, error) {
	return m.GenerateAccessTokenWithClaims(userID, sessionID, uuid.Nil, nil)
}

// GenerateAccessTokenWithClaims generates a JWT access token with custom claims.
func (m *Manager) GenerateAccessTokenWithClaims(
	userID uuid.UUID,
	sessionID uuid.UUID,
	orgID uuid.UUID,
	claims *auth.TokenClaims,
) (string, error) {
	now := time.Now()
	exp := now.Add(1 * time.Hour)

	// Create enhanced claims
	enhancedClaims := &auth.EnhancedClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:         userID,
		TokenType:      auth.AccessTokenType,
		SessionID:      sessionID.String(),
		EmailVerified:  true,
		AuthMethod:     auth.AuthMethodOAuth,
		OrganizationID: orgID,
	}

	// Add user context if available
	if claims != nil {
		enhancedClaims.Email = claims.Email
		enhancedClaims.Name = claims.Name
		enhancedClaims.AvatarURL = claims.Picture
		enhancedClaims.Provider = claims.Provider
		enhancedClaims.ProviderID = claims.ProviderID
		enhancedClaims.OrganizationName = claims.OrganizationName
		enhancedClaims.OrganizationRole = claims.OrganizationRole
		enhancedClaims.Roles = claims.Roles
		enhancedClaims.Permissions = claims.Permissions
		enhancedClaims.Scopes = claims.Scopes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, enhancedClaims)
	tokenString, err := token.SignedString([]byte(m.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a JWT refresh token.
func (m *Manager) GenerateRefreshToken(userID uuid.UUID, sessionID uuid.UUID) (string, error) {
	now := time.Now()
	exp := now.Add(30 * 24 * time.Hour) // 30 days

	claims := &auth.RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:     userID,
		TokenType:  auth.RefreshTokenType,
		SessionID:  sessionID.String(),
		AuthMethod: auth.AuthMethodOAuth,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateEnhancedToken validates a JWT token and returns the enhanced claims.
func (m *Manager) ValidateEnhancedToken(tokenString string) (*auth.EnhancedClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&auth.EnhancedClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*auth.EnhancedClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate claims
	if !claims.IsValid() {
		return nil, fmt.Errorf("invalid or expired claims")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token and returns the claims.
func (m *Manager) ValidateRefreshToken(tokenString string) (*auth.RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&auth.RefreshClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*auth.RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	// Check expiration
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("refresh token expired")
	}

	return claims, nil
}

// GenerateAPIToken generates an API token (not implemented - use APITokenStore).
func (m *Manager) GenerateAPIToken(_ uuid.UUID, _ string, _ []string) (string, error) {
	// API tokens are managed by APITokenStore, not JWT
	return "", fmt.Errorf("API token generation should use APITokenStore")
}

// RefreshToken handles token refresh (interface method).
func (m *Manager) RefreshToken(refreshToken string) (string, string, error) {
	// Validate refresh token
	refreshClaims, err := m.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("validate refresh token: %w", err)
	}

	// Parse user ID
	userID, err := uuid.Parse(refreshClaims.UserID.String())
	if err != nil {
		return "", "", fmt.Errorf("parse user ID: %w", err)
	}

	// Parse session ID
	sessionID, err := uuid.Parse(refreshClaims.SessionID)
	if err != nil {
		return "", "", fmt.Errorf("parse session ID: %w", err)
	}

	// Generate new token pair
	newAccessToken, err := m.GenerateAccessToken(userID, sessionID)
	if err != nil {
		return "", "", fmt.Errorf("generate new access token: %w", err)
	}

	newRefreshToken, err := m.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		return "", "", fmt.Errorf("generate new refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// ValidateToken validates a JWT token and returns auth.TokenClaims (interface method).
func (m *Manager) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	enhancedClaims, err := m.ValidateEnhancedToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Convert to auth.TokenClaims
	return &auth.TokenClaims{
		UserID:           enhancedClaims.UserID,
		SessionID:        uuid.MustParse(enhancedClaims.SessionID),
		Email:            enhancedClaims.Email,
		Name:             enhancedClaims.Name,
		Picture:          enhancedClaims.AvatarURL,
		Provider:         enhancedClaims.Provider,
		ProviderID:       enhancedClaims.ProviderID,
		OrganizationName: enhancedClaims.OrganizationName,
		OrganizationRole: enhancedClaims.OrganizationRole,
		Roles:            enhancedClaims.Roles,
		Permissions:      enhancedClaims.Permissions,
		Scopes:           enhancedClaims.Scopes,
		ExpiresAt:        enhancedClaims.ExpiresAt.Time,
	}, nil
}

// RefreshTokenPair generates a new token pair using a refresh token.
func (m *Manager) RefreshTokenPair(
	refreshToken string,
	sessionID uuid.UUID,
	claims *auth.TokenClaims,
) (*auth.TokenPair, error) {
	// Validate refresh token
	refreshClaims, err := m.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("validate refresh token: %w", err)
	}

	// Verify session ID matches
	if refreshClaims.SessionID != sessionID.String() {
		return nil, fmt.Errorf("session ID mismatch")
	}

	// Generate new access token
	userID, err := uuid.Parse(refreshClaims.Subject)
	if err != nil {
		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	accessToken, err := m.GenerateAccessTokenWithClaims(userID, sessionID, uuid.Nil, claims)
	if err != nil {
		return nil, fmt.Errorf("generate new access token: %w", err)
	}

	return &auth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // Reuse the same refresh token
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// GetUserIDFromToken extracts the user ID from a token without full validation.
func (m *Manager) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims")
	}

	userIDStr, ok := claims["uid"].(string)
	if !ok {
		// Try "sub" as fallback
		userIDStr, ok = claims["sub"].(string)
		if !ok {
			return uuid.Nil, fmt.Errorf("user ID not found in claims")
		}
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse user ID: %w", err)
	}

	return userID, nil
}
