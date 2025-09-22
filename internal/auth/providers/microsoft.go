package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/archesai/archesai/internal/auth"
)

// MicrosoftProvider implements OAuth for Microsoft
type MicrosoftProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

// NewMicrosoftProvider creates a new Microsoft OAuth provider
func NewMicrosoftProvider(clientID, clientSecret, redirectURL string) *MicrosoftProvider {
	return &MicrosoftProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

// GetAuthURL returns the Microsoft OAuth authorization URL
func (p *MicrosoftProvider) GetAuthURL(state string) string {
	params := url.Values{
		"client_id":     {p.clientID},
		"redirect_uri":  {p.redirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile offline_access"},
		"state":         {state},
		"response_mode": {"query"},
	}
	return fmt.Sprintf(
		"https://login.microsoftonline.com/common/oauth2/v2.0/authorize?%s",
		params.Encode(),
	)
}

// ExchangeCode exchanges authorization code for tokens
func (p *MicrosoftProvider) ExchangeCode(
	ctx context.Context,
	code string,
) (*auth.OAuthTokens, error) {
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
		"https://login.microsoftonline.com/common/oauth2/v2.0/token",
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

	return &auth.OAuthTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		IDToken:      tokenResp.IDToken,
		ExpiresIn:    tokenResp.ExpiresIn,
		Scope:        tokenResp.Scope,
	}, nil
}

// GetUserInfo fetches user information from Microsoft Graph
func (p *MicrosoftProvider) GetUserInfo(
	ctx context.Context,
	accessToken string,
) (*auth.OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://graph.microsoft.com/v1.0/me",
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
		ID                string `json:"id"`
		DisplayName       string `json:"displayName"`
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	// Microsoft might return email in mail or userPrincipalName
	email := userInfo.Mail
	if email == "" {
		email = userInfo.UserPrincipalName
	}

	return &auth.OAuthUserInfo{
		ID:            userInfo.ID,
		Email:         email,
		Name:          userInfo.DisplayName,
		Picture:       "",   // Microsoft Graph requires separate endpoint for profile photo
		EmailVerified: true, // Microsoft accounts are verified
		Provider:      auth.ProviderMicrosoft,
	}, nil
}

// RefreshToken refreshes an expired access token
func (p *MicrosoftProvider) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*auth.OAuthTokens, error) {
	data := url.Values{
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://login.microsoftonline.com/common/oauth2/v2.0/token",
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

	return &auth.OAuthTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: refreshToken, // Reuse the same refresh token
		ExpiresIn:    tokenResp.ExpiresIn,
		Scope:        tokenResp.Scope,
	}, nil
}
