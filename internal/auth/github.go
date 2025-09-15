// Package auth provides OAuth2 provider implementations for various services.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	githubUserURL = "https://api.github.com/user"
)

// GitHubOAuthProvider implements OAuth2 for GitHub
type GitHubOAuthProvider struct {
	*BaseOAuthProvider
}

// NewGitHubOAuthProvider creates a new GitHub OAuth provider
func NewGitHubOAuthProvider(clientID, clientSecret string) OAuthProvider {
	return &GitHubOAuthProvider{
		BaseOAuthProvider: &BaseOAuthProvider{
			Config: &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     github.Endpoint,
				Scopes: []string{
					"user:email",
					"read:user",
				},
			},
		},
	}
}

// GetProviderID returns the provider identifier
func (p *GitHubOAuthProvider) GetProviderID() string {
	return "github"
}

// GetAuthURL returns the GitHub authorization URL
func (p *GitHubOAuthProvider) GetAuthURL(state string, redirectURI string) string {
	p.Config.RedirectURL = redirectURI
	return p.Config.AuthCodeURL(state)
}

// ExchangeCode exchanges an authorization code for tokens
func (p *GitHubOAuthProvider) ExchangeCode(ctx context.Context, code string, redirectURI string) (*OAuthTokens, error) {
	p.Config.RedirectURL = redirectURI
	token, err := p.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &OAuthTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken, // GitHub doesn't provide refresh tokens
		ExpiresIn:    int(token.Expiry.Unix()),
		Scope:        token.Extra("scope").(string),
	}, nil
}

// RefreshToken is not supported by GitHub
func (p *GitHubOAuthProvider) RefreshToken(_ context.Context, _ string) (*OAuthTokens, error) {
	// GitHub doesn't support refresh tokens
	return nil, fmt.Errorf("GitHub does not support refresh tokens")
}

// GetUserInfo retrieves user information from GitHub
func (p *GitHubOAuthProvider) GetUserInfo(ctx context.Context, accessToken string) (*OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", githubUserURL, nil)
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
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var githubUser struct {
		ID                int    `json:"id"`
		Login             string `json:"login"`
		Email             string `json:"email"`
		Name              string `json:"name"`
		AvatarURL         string `json:"avatar_url"`
		Location          string `json:"location"`
		PublicRepos       int    `json:"public_repos"`
		PublicGists       int    `json:"public_gists"`
		Followers         int    `json:"followers"`
		Following         int    `json:"following"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at"`
		PrivateGists      int    `json:"private_gists"`
		TotalPrivateRepos int    `json:"total_private_repos"`
		OwnedPrivateRepos int    `json:"owned_private_repos"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// GitHub may not return email in the user endpoint if it's private
	// In production, you might want to make an additional call to /user/emails
	email := githubUser.Email
	if email == "" {
		// Get primary email from emails endpoint
		email, err = p.getPrimaryEmail(ctx, accessToken)
		if err != nil {
			// Non-fatal: user might not have granted email permission
			email = fmt.Sprintf("%s@users.noreply.github.com", githubUser.Login)
		}
	}

	return &OAuthUserInfo{
		ProviderAccountID: fmt.Sprintf("%d", githubUser.ID),
		Email:             email,
		EmailVerified:     true, // GitHub requires email verification
		Name:              githubUser.Name,
		Picture:           githubUser.AvatarURL,
		Raw: map[string]interface{}{
			"id":         githubUser.ID,
			"login":      githubUser.Login,
			"email":      email,
			"name":       githubUser.Name,
			"avatar_url": githubUser.AvatarURL,
			"location":   githubUser.Location,
			"created_at": githubUser.CreatedAt,
			"updated_at": githubUser.UpdatedAt,
		},
	}, nil
}

// getPrimaryEmail retrieves the primary email from GitHub's emails endpoint
func (p *GitHubOAuthProvider) getPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get emails: status %d", resp.StatusCode)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

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
