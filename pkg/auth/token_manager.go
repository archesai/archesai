package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenManager handles JWT token creation and validation.
type TokenManager struct {
	jwtSecret []byte
}

// NewTokenManager creates a new token manager.
func NewTokenManager(jwtSecret string) *TokenManager {
	return &TokenManager{
		jwtSecret: []byte(jwtSecret),
	}
}

// TokenClaims represents the claims in a JWT token.
type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

// CreateAccessToken creates a short-lived access token.
func (tm *TokenManager) CreateAccessToken(userID, sessionID uuid.UUID) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.jwtSecret)
}

// CreateRefreshToken creates a long-lived refresh token.
func (tm *TokenManager) CreateRefreshToken(userID, sessionID uuid.UUID) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.jwtSecret)
}

// ValidateAccessToken validates an access token and returns claims.
func (tm *TokenManager) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	return tm.validateToken(tokenString, "access")
}

// ValidateRefreshToken validates a refresh token and returns claims.
func (tm *TokenManager) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	return tm.validateToken(tokenString, "refresh")
}

// validateToken validates a token and checks its type.
func (tm *TokenManager) validateToken(
	tokenString string,
	expectedType string,
) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return tm.jwtSecret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if claims.TokenType != expectedType {
			return nil, fmt.Errorf(
				"invalid token type: expected %s, got %s",
				expectedType,
				claims.TokenType,
			)
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
