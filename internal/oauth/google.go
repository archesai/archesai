package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleProvider implements OAuth2 for Google
type GoogleProvider struct {
	*BaseOAuthProvider
	providerID string
	logger     *slog.Logger
}

// NewGoogleProvider creates a new Google OAuth provider
func NewGoogleProvider(
	clientID, clientSecret, redirectURL string,
	logger *slog.Logger,
) *GoogleProvider {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleProvider{
		BaseOAuthProvider: &BaseOAuthProvider{
			Config: config,
		},
		providerID: "google",
		logger:     logger,
	}
}

// GetAuthURL returns the authorization URL for the provider
func (p *GoogleProvider) GetAuthURL(state string, redirectURI string) string {
	return p.BaseOAuthProvider.GetAuthURL(state, redirectURI)
}

// ExchangeCode exchanges an authorization code for tokens
func (p *GoogleProvider) ExchangeCode(
	ctx context.Context,
	code string,
	redirectURI string,
) (*Tokens, error) {
	token, err := p.BaseOAuthProvider.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	idToken := ""
	if idTokenVal := token.Extra("id_token"); idTokenVal != nil {
		idToken, _ = idTokenVal.(string)
	}

	return &Tokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      idToken,
		ExpiresIn:    int(time.Until(token.Expiry).Seconds()),
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (p *GoogleProvider) RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error) {
	token, err := p.BaseOAuthProvider.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &Tokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(time.Until(token.Expiry).Seconds()),
	}, nil
}

// GetUserInfo retrieves user information from Google
func (p *GoogleProvider) GetUserInfo(_ context.Context, accessToken string) (*UserInfo, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accessToken),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Error("failed to close response body", "error", err)
		}
	}()

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &UserInfo{
		ProviderAccountID: googleUser.ID,
		Email:             googleUser.Email,
		EmailVerified:     googleUser.VerifiedEmail,
		Name:              googleUser.Name,
		Picture:           googleUser.Picture,
		Locale:            googleUser.Locale,
	}, nil
}

// GetProviderID returns the provider identifier
func (p *GoogleProvider) GetProviderID() string {
	return p.providerID
}
