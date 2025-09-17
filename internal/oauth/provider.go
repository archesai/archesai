package oauth

import (
	"context"
	"net/url"

	"golang.org/x/oauth2"
)

// Provider defines the interface for OAuth2 providers.
type Provider interface {
	// GetAuthURL returns the authorization URL for the provider
	GetAuthURL(state string, redirectURI string) string

	// ExchangeCode exchanges an authorization code for tokens
	ExchangeCode(ctx context.Context, code string, redirectURI string) (*Tokens, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error)

	// GetUserInfo retrieves user information from the provider
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)

	// GetProviderID returns the provider identifier (e.g., "google", "github")
	GetProviderID() string
}

// Tokens represents OAuth2 tokens returned by a provider.
type Tokens struct {
	AccessToken  string
	RefreshToken string
	IDToken      string // OpenID Connect ID token (optional)
	ExpiresIn    int    // Seconds until expiration
	Scope        string
}

// UserInfo represents user information from an OAuth provider.
type UserInfo struct {
	ProviderAccountID string // Unique ID from the provider
	Email             string
	EmailVerified     bool
	Name              string
	Picture           string
	Locale            string
	// Provider-specific additional data
	Raw map[string]interface{}
}

// Config holds common OAuth2 configuration.
type Config struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	AuthURL      string
	TokenURL     string
}

// BaseOAuthProvider provides common OAuth2 functionality.
type BaseOAuthProvider struct {
	Config *oauth2.Config
}

// GetAuthURL returns the authorization URL with the provided state.
func (p *BaseOAuthProvider) GetAuthURL(state string, redirectURI string) string {
	p.Config.RedirectURL = redirectURI

	// Add access_type=offline for providers that support refresh tokens
	authURL, _ := url.Parse(p.Config.AuthCodeURL(state))
	q := authURL.Query()
	q.Set("access_type", "offline")
	q.Set("prompt", "consent") // Force consent to get refresh token
	authURL.RawQuery = q.Encode()

	return authURL.String()
}

// ExchangeCode exchanges an authorization code for tokens.
func (p *BaseOAuthProvider) ExchangeCode(
	ctx context.Context,
	code string,
	redirectURI string,
) (*oauth2.Token, error) {
	p.Config.RedirectURL = redirectURI
	return p.Config.Exchange(ctx, code)
}

// RefreshToken refreshes an access token using a refresh token.
func (p *BaseOAuthProvider) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*oauth2.Token, error) {
	tokenSource := p.Config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})
	return tokenSource.Token()
}
