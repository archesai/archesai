// Package oauth implements OAuth providers for authentication
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GitHubProvider implements OAuth for GitHub
type GitHubProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

// NewGitHubProvider creates a new GitHub OAuth provider
func NewGitHubProvider(clientID, clientSecret, redirectURL string) *GitHubProvider {
	return &GitHubProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

// GetAuthURL returns the GitHub OAuth authorization URL
func (p *GitHubProvider) GetAuthURL(state string) string {
	params := url.Values{
		"client_id":    {p.clientID},
		"redirect_uri": {p.redirectURL},
		"scope":        {"user:email"},
		"state":        {state},
	}
	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params.Encode())
}

// ExchangeCode exchanges authorization code for tokens
func (p *GitHubProvider) ExchangeCode(
	ctx context.Context,
	code string,
) (*Tokens, error) {
	// Prepare token exchange request
	data := url.Values{
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
		"code":          {code},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://github.com/login/oauth/access_token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s", body)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: "", // GitHub doesn't use refresh tokens
		IDToken:      "", // GitHub doesn't provide ID tokens
		ExpiresIn:    0,  // GitHub tokens don't expire
		Scope:        tokenResp.Scope,
	}, nil
}

// GetUserInfo fetches user information from GitHub
func (p *GitHubProvider) GetUserInfo(
	ctx context.Context,
	accessToken string,
) (*UserInfo, error) {
	// Get user info
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user info failed: %s", body)
	}

	var userInfo struct {
		ID     int    `json:"id"`
		Login  string `json:"login"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar_url"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	// GitHub may not return email in the user endpoint, fetch from emails endpoint
	email := userInfo.Email
	if email == "" {
		email, err = p.getPrimaryEmail(ctx, accessToken)
		if err != nil {
			// Continue without email if we can't fetch it
			email = fmt.Sprintf("%s@users.noreply.github.com", userInfo.Login)
		}
	}

	// Use login as name if name is empty
	name := userInfo.Name
	if name == "" {
		name = userInfo.Login
	}

	return &UserInfo{
		ID:            fmt.Sprintf("%d", userInfo.ID),
		Email:         email,
		Name:          name,
		Picture:       userInfo.Avatar,
		EmailVerified: true, // GitHub requires email verification
		Provider:      "github",
	}, nil
}

// RefreshToken refreshes an expired access token (GitHub doesn't support this)
func (p *GitHubProvider) RefreshToken(
	_ context.Context,
	_ string,
) (*Tokens, error) {
	return nil, fmt.Errorf("GitHub does not support token refresh")
}

// getPrimaryEmail fetches the primary email from GitHub
func (p *GitHubProvider) getPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get emails failed: %s", body)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	// Find primary email
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	// Fallback to any verified email
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}
