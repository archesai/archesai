package auth

import (
	"context"
	"net/url"

	"golang.org/x/oauth2"
)

// OAuthProvider defines the interface for OAuth2 providers
type OAuthProvider interface {
	// GetAuthURL returns the authorization URL for the provider
	GetAuthURL(state string, redirectURI string) string

	// ExchangeCode exchanges an authorization code for tokens
	ExchangeCode(ctx context.Context, code string, redirectURI string) (*OAuthTokens, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*OAuthTokens, error)

	// GetUserInfo retrieves user information from the provider
	GetUserInfo(ctx context.Context, accessToken string) (*OAuthUserInfo, error)

	// GetProviderID returns the provider identifier (e.g., "google", "github")
	GetProviderID() string
}

// OAuthTokens represents OAuth2 tokens returned by a provider
type OAuthTokens struct {
	AccessToken  string
	RefreshToken string
	IDToken      string // OpenID Connect ID token (optional)
	ExpiresIn    int    // Seconds until expiration
	Scope        string
}

// OAuthUserInfo represents user information from an OAuth provider
type OAuthUserInfo struct {
	ProviderAccountID string // Unique ID from the provider
	Email             string
	EmailVerified     bool
	Name              string
	Picture           string
	Locale            string
	// Provider-specific additional data
	Raw map[string]interface{}
}

// OAuthConfig holds common OAuth2 configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	AuthURL      string
	TokenURL     string
}

// baseOAuthProvider provides common OAuth2 functionality
type baseOAuthProvider struct {
	config *oauth2.Config
}

// GetAuthURL returns the authorization URL with the provided state
func (p *baseOAuthProvider) GetAuthURL(state string, redirectURI string) string {
	p.config.RedirectURL = redirectURI

	// Add access_type=offline for providers that support refresh tokens
	authURL, _ := url.Parse(p.config.AuthCodeURL(state))
	q := authURL.Query()
	q.Set("access_type", "offline")
	q.Set("prompt", "consent") // Force consent to get refresh token
	authURL.RawQuery = q.Encode()

	return authURL.String()
}

// ExchangeCode exchanges an authorization code for tokens
func (p *baseOAuthProvider) ExchangeCode(ctx context.Context, code string, redirectURI string) (*oauth2.Token, error) {
	p.config.RedirectURL = redirectURI
	return p.config.Exchange(ctx, code)
}

// RefreshToken refreshes an access token using a refresh token
func (p *baseOAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	tokenSource := p.config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})
	return tokenSource.Token()
}
