package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/archesai/archesai/internal/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

const (
	microsoftUserInfoURL = "https://graph.microsoft.com/v1.0/me"
)

// MicrosoftOAuthProvider implements OAuth2 for Microsoft
type MicrosoftOAuthProvider struct {
	*auth.BaseOAuthProvider
}

// NewMicrosoftOAuthProvider creates a new Microsoft OAuth provider
func NewMicrosoftOAuthProvider(clientID, clientSecret string) auth.OAuthProvider {
	return &MicrosoftOAuthProvider{
		BaseOAuthProvider: &auth.BaseOAuthProvider{
			Config: &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     microsoft.AzureADEndpoint("common"), // "common" allows any Azure AD and personal Microsoft account
				Scopes: []string{
					"openid",
					"profile",
					"email",
					"offline_access", // For refresh token
					"User.Read",      // Microsoft Graph API permission
				},
			},
		},
	}
}

// GetProviderID returns the provider identifier
func (p *MicrosoftOAuthProvider) GetProviderID() string {
	return "microsoft"
}

// GetAuthURL returns the Microsoft authorization URL
func (p *MicrosoftOAuthProvider) GetAuthURL(state string, redirectURI string) string {
	p.Config.RedirectURL = redirectURI
	// Add prompt=select_account to allow account selection
	return p.Config.AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "select_account"))
}

// ExchangeCode exchanges an authorization code for tokens
func (p *MicrosoftOAuthProvider) ExchangeCode(ctx context.Context, code string, redirectURI string) (*auth.OAuthTokens, error) {
	p.Config.RedirectURL = redirectURI
	token, err := p.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &auth.OAuthTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      extractIDToken(token),
		ExpiresIn:    int(token.Expiry.Unix()),
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (p *MicrosoftOAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*auth.OAuthTokens, error) {
	token, err := p.BaseOAuthProvider.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &auth.OAuthTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.Expiry.Unix()),
	}, nil
}

// GetUserInfo retrieves user information from Microsoft Graph
func (p *MicrosoftOAuthProvider) GetUserInfo(ctx context.Context, accessToken string) (*auth.OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", microsoftUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

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

	var msUser struct {
		ID                string   `json:"id"`
		DisplayName       string   `json:"displayName"`
		GivenName         string   `json:"givenName"`
		Surname           string   `json:"surname"`
		Mail              string   `json:"mail"`
		UserPrincipalName string   `json:"userPrincipalName"`
		PreferredLanguage string   `json:"preferredLanguage"`
		MobilePhone       string   `json:"mobilePhone"`
		JobTitle          string   `json:"jobTitle"`
		OfficeLocation    string   `json:"officeLocation"`
		BusinessPhones    []string `json:"businessPhones"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&msUser); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Microsoft might return mail or userPrincipalName
	email := msUser.Mail
	if email == "" {
		email = msUser.UserPrincipalName
	}

	// Get profile photo URL (optional, may fail for some accounts)
	photoURL, _ := p.getProfilePhotoURL(ctx, accessToken)

	return &auth.OAuthUserInfo{
		ProviderAccountID: msUser.ID,
		Email:             email,
		EmailVerified:     true, // Microsoft requires email verification
		Name:              msUser.DisplayName,
		Picture:           photoURL,
		Locale:            msUser.PreferredLanguage,
		Raw: map[string]interface{}{
			"id":                msUser.ID,
			"displayName":       msUser.DisplayName,
			"givenName":         msUser.GivenName,
			"surname":           msUser.Surname,
			"mail":              msUser.Mail,
			"userPrincipalName": msUser.UserPrincipalName,
			"preferredLanguage": msUser.PreferredLanguage,
			"mobilePhone":       msUser.MobilePhone,
			"jobTitle":          msUser.JobTitle,
			"officeLocation":    msUser.OfficeLocation,
			"businessPhones":    msUser.BusinessPhones,
		},
	}, nil
}

// getProfilePhotoURL attempts to get the user's profile photo URL
func (p *MicrosoftOAuthProvider) getProfilePhotoURL(ctx context.Context, accessToken string) (string, error) {
	// Check if photo exists
	req, err := http.NewRequestWithContext(ctx, "GET", "https://graph.microsoft.com/v1.0/me/photo", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		// Photo exists, return the URL to fetch it
		// Note: In production, you might want to actually fetch and store the photo
		return "https://graph.microsoft.com/v1.0/me/photo/$value", nil
	}

	return "", fmt.Errorf("no profile photo available")
}
