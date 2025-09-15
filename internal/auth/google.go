package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// GoogleOAuthProvider implements OAuth2 for Google
type GoogleOAuthProvider struct {
	*BaseOAuthProvider
}

// Ensure GoogleOAuthProvider implements OAuthProvider
var _ OAuthProvider = (*GoogleOAuthProvider)(nil)

// NewGoogleOAuthProvider creates a new Google OAuth provider
func NewGoogleOAuthProvider(clientID, clientSecret string) OAuthProvider {
	return &GoogleOAuthProvider{
		BaseOAuthProvider: &BaseOAuthProvider{
			Config: &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     google.Endpoint,
				Scopes: []string{
					"openid",
					"profile",
					"email",
				},
			},
		},
	}
}

// GetProviderID returns the provider identifier
func (p *GoogleOAuthProvider) GetProviderID() string {
	return "google"
}

// GetAuthURL returns the Google authorization URL
func (p *GoogleOAuthProvider) GetAuthURL(state string, redirectURI string) string {
	p.Config.RedirectURL = redirectURI
	// Google requires access_type=offline to get refresh token
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// ExchangeCode exchanges an authorization code for tokens
func (p *GoogleOAuthProvider) ExchangeCode(ctx context.Context, code string, redirectURI string) (*OAuthTokens, error) {
	p.Config.RedirectURL = redirectURI
	token, err := p.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &OAuthTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      extractIDToken(token),
		ExpiresIn:    int(token.Expiry.Unix()),
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (p *GoogleOAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*OAuthTokens, error) {
	token, err := p.BaseOAuthProvider.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &OAuthTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.Expiry.Unix()),
	}, nil
}

// GetUserInfo retrieves user information from Google
func (p *GoogleOAuthProvider) GetUserInfo(ctx context.Context, accessToken string) (*OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", googleUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &OAuthUserInfo{
		ProviderAccountID: googleUser.ID,
		Email:             googleUser.Email,
		EmailVerified:     googleUser.VerifiedEmail,
		Name:              googleUser.Name,
		Picture:           googleUser.Picture,
		Locale:            googleUser.Locale,
		Raw: map[string]interface{}{
			"id":             googleUser.ID,
			"email":          googleUser.Email,
			"verified_email": googleUser.VerifiedEmail,
			"name":           googleUser.Name,
			"picture":        googleUser.Picture,
			"locale":         googleUser.Locale,
		},
	}, nil
}

// extractIDToken extracts the ID token from the OAuth2 token extra fields
func extractIDToken(token *oauth2.Token) string {
	if idToken, ok := token.Extra("id_token").(string); ok {
		return idToken
	}
	return ""
}
