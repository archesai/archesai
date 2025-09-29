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

// GoogleProvider implements OAuth for Google
type GoogleProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

// NewGoogleProvider creates a new Google OAuth provider
func NewGoogleProvider(clientID, clientSecret, redirectURL string) *GoogleProvider {
	return &GoogleProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

// GetAuthURL returns the Google OAuth authorization URL
func (p *GoogleProvider) GetAuthURL(state string) string {
	params := url.Values{
		"client_id":     {p.clientID},
		"redirect_uri":  {p.redirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
	}
	return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?%s", params.Encode())
}

// ExchangeCode exchanges authorization code for tokens
func (p *GoogleProvider) ExchangeCode(
	ctx context.Context,
	code string,
) (*OAuthTokens, error) {
	// Prepare token exchange request
	data := url.Values{
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {p.redirectURL},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
		IDToken      string `json:"id_token,omitempty"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &OAuthTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		IDToken:      tokenResp.IDToken,
		ExpiresIn:    tokenResp.ExpiresIn,
		Scope:        tokenResp.Scope,
	}, nil
}

// GetUserInfo fetches user information from Google
func (p *GoogleProvider) GetUserInfo(
	ctx context.Context,
	accessToken string,
) (*OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://www.googleapis.com/oauth2/v2/userinfo",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		EmailVerified bool   `json:"verified_email"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &OAuthUserInfo{
		ID:            userInfo.ID,
		Email:         userInfo.Email,
		Name:          userInfo.Name,
		Picture:       userInfo.Picture,
		EmailVerified: userInfo.EmailVerified,
		Provider:      "google",
	}, nil
}

// RefreshToken refreshes an expired access token
func (p *GoogleProvider) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*OAuthTokens, error) {
	data := url.Values{
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		return nil, fmt.Errorf("token refresh failed: %s", body)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &OAuthTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: refreshToken, // Reuse the same refresh token
		ExpiresIn:    tokenResp.ExpiresIn,
		Scope:        tokenResp.Scope,
	}, nil
}
