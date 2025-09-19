// Package oauth provides OAuth2 authentication support for multiple providers.
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubProvider implements OAuth2 for GitHub
type GitHubProvider struct {
	*BaseOAuthProvider
	providerID string
	logger     *slog.Logger
}

// NewGitHubProvider creates a new GitHub OAuth provider
func NewGitHubProvider(
	clientID, clientSecret, redirectURL string,
	logger *slog.Logger,
) *GitHubProvider {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"read:user",
			"user:email",
		},
		Endpoint: github.Endpoint,
	}

	return &GitHubProvider{
		BaseOAuthProvider: &BaseOAuthProvider{
			Config: config,
		},
		providerID: "github",
		logger:     logger,
	}
}

// GetAuthURL returns the authorization URL for GitHub
func (p *GitHubProvider) GetAuthURL(state, redirectURI string) string {
	return p.BaseOAuthProvider.GetAuthURL(state, redirectURI)
}

// ExchangeCode exchanges an authorization code for tokens
func (p *GitHubProvider) ExchangeCode(
	ctx context.Context,
	code string,
	redirectURI string,
) (*Tokens, error) {
	token, err := p.BaseOAuthProvider.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &Tokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(time.Until(token.Expiry).Seconds()),
	}, nil
}

// GetUserInfo retrieves user information from GitHub
func (p *GitHubProvider) GetUserInfo(_ context.Context, accessToken string) (*UserInfo, error) {
	// Get user info
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Error("failed to close response body", "error", err)
		}
	}()

	var ghUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Location  string `json:"location"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// If email is not public, we need to fetch it separately
	if ghUser.Email == "" {
		email, err := p.getUserEmail(accessToken)
		if err == nil {
			ghUser.Email = email
		}
	}

	// Use login as name if name is not set
	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}

	return &UserInfo{
		ProviderAccountID: fmt.Sprintf("%d", ghUser.ID),
		Email:             ghUser.Email,
		EmailVerified:     true, // GitHub verifies emails
		Name:              name,
		Picture:           ghUser.AvatarURL,
	}, nil
}

// getUserEmail fetches the primary email from GitHub
func (p *GitHubProvider) getUserEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get user emails: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.logger.Error("failed to close response body", "error", err)
		}
	}()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode emails: %w", err)
	}

	// Find primary email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// Return first verified email if no primary
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}

// RefreshToken refreshes an access token using a refresh token
func (p *GitHubProvider) RefreshToken(_ context.Context, _ string) (*Tokens, error) {
	// GitHub doesn't support refresh tokens in OAuth apps
	// You need to use GitHub Apps for refresh tokens
	return nil, fmt.Errorf("GitHub OAuth apps don't support refresh tokens")
}

// GetProviderID returns the provider identifier
func (p *GitHubProvider) GetProviderID() string {
	return p.providerID
}
