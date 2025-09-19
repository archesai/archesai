package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

// MicrosoftProvider implements OAuth2 for Microsoft
type MicrosoftProvider struct {
	*BaseOAuthProvider
	providerID string
	logger     *slog.Logger
}

// NewMicrosoftProvider creates a new Microsoft OAuth provider
func NewMicrosoftProvider(
	clientID, clientSecret, redirectURL string,
	logger *slog.Logger,
) *MicrosoftProvider {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"openid",
			"email",
			"profile",
			"offline_access",
		},
		Endpoint: microsoft.AzureADEndpoint("common"),
	}

	return &MicrosoftProvider{
		BaseOAuthProvider: &BaseOAuthProvider{
			Config: config,
		},
		providerID: "microsoft",
		logger:     logger,
	}
}

// GetAuthURL returns the authorization URL for Microsoft
func (p *MicrosoftProvider) GetAuthURL(state, redirectURI string) string {
	return p.BaseOAuthProvider.GetAuthURL(state, redirectURI)
}

// ExchangeCode exchanges an authorization code for tokens
func (p *MicrosoftProvider) ExchangeCode(
	ctx context.Context,
	code string,
	redirectURI string,
) (*Tokens, error) {
	token, err := p.BaseOAuthProvider.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	idToken := ""
	if idTokenRaw := token.Extra("id_token"); idTokenRaw != nil {
		idToken = idTokenRaw.(string)
	}

	return &Tokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      idToken,
		ExpiresIn:    int(time.Until(token.Expiry).Seconds()),
	}, nil
}

// GetUserInfo retrieves user information from Microsoft Graph
func (p *MicrosoftProvider) GetUserInfo(
	_ context.Context,
	accessToken string,
) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
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
		if err := resp.Body.Close(); err != nil {
			p.logger.Error("failed to close response body", "error", err)
		}
	}()

	var msUser struct {
		ID                string `json:"id"`
		DisplayName       string `json:"displayName"`
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&msUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Microsoft might return email in mail or userPrincipalName
	email := msUser.Mail
	if email == "" {
		email = msUser.UserPrincipalName
	}

	return &UserInfo{
		ProviderAccountID: msUser.ID,
		Email:             email,
		EmailVerified:     true, // Microsoft accounts are verified
		Name:              msUser.DisplayName,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (p *MicrosoftProvider) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*Tokens, error) {
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

// GetProviderID returns the provider identifier
func (p *MicrosoftProvider) GetProviderID() string {
	return p.providerID
}
