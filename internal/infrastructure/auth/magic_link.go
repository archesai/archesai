// Package auth provides authentication and authorization implementations
package auth

import (
	"fmt"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MagicLinkProvider handles magic link generation and validation.
type MagicLinkProvider struct {
	jwtSecret []byte
	baseURL   string
}

// NewMagicLinkProvider creates a new magic link provider.
func NewMagicLinkProvider(jwtSecret, baseURL string) *MagicLinkProvider {
	return &MagicLinkProvider{
		jwtSecret: []byte(jwtSecret),
		baseURL:   baseURL,
	}
}

// MagicLinkClaims represents the claims in a magic link token.
type MagicLinkClaims struct {
	Identifier  string `json:"identifier"`
	Purpose     string `json:"purpose"`
	RedirectURL string `json:"redirect_url,omitempty"`
	jwt.RegisteredClaims
}

// GenerateLink creates a stateless magic link token.
func (mlp *MagicLinkProvider) GenerateLink(identifier, redirectURL string) (string, error) {
	claims := MagicLinkClaims{
		Identifier:  identifier,
		Purpose:     "login",
		RedirectURL: redirectURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mlp.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign magic link token: %w", err)
	}

	// Build the magic link URL
	u, err := url.Parse(mlp.baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = "/auth/magic-links/verify"
	q := u.Query()
	q.Set("token", tokenString)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// ValidateLink validates a magic link token.
func (mlp *MagicLinkProvider) ValidateLink(tokenString string) (*MagicLinkClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MagicLinkClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return mlp.jwtSecret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse magic link token: %w", err)
	}

	if claims, ok := token.Claims.(*MagicLinkClaims); ok && token.Valid {
		// Check if token is expired
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, fmt.Errorf("magic link has expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid magic link token")
}
