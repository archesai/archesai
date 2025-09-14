package oauth

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestGoogleOAuthProvider_GetProviderID(t *testing.T) {
	provider := NewGoogleOAuthProvider("test-client-id", "test-client-secret")
	assert.Equal(t, "google", provider.GetProviderID())
}

func TestGoogleOAuthProvider_GetAuthURL(t *testing.T) {
	provider := NewGoogleOAuthProvider("test-client-id", "test-client-secret")

	state := "test-state-123"
	redirectURI := "http://localhost:8080/auth/callback/google"

	authURL := provider.GetAuthURL(state, redirectURI)

	// Parse the URL
	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	// Check base URL
	assert.Equal(t, "https", parsedURL.Scheme)
	assert.Equal(t, "accounts.google.com", parsedURL.Host)
	assert.Equal(t, "/o/oauth2/auth", parsedURL.Path)

	// Check query parameters
	query := parsedURL.Query()
	assert.Equal(t, "test-client-id", query.Get("client_id"))
	assert.Equal(t, redirectURI, query.Get("redirect_uri"))
	assert.Equal(t, "code", query.Get("response_type"))
	assert.Equal(t, state, query.Get("state"))
	assert.Equal(t, "offline", query.Get("access_type"))
	assert.Contains(t, query.Get("scope"), "openid")
	assert.Contains(t, query.Get("scope"), "email")
	assert.Contains(t, query.Get("scope"), "profile")
}

func TestGitHubOAuthProvider_GetProviderID(t *testing.T) {
	provider := NewGitHubOAuthProvider("test-client-id", "test-client-secret")
	assert.Equal(t, "github", provider.GetProviderID())
}

func TestGitHubOAuthProvider_GetAuthURL(t *testing.T) {
	provider := NewGitHubOAuthProvider("test-client-id", "test-client-secret")

	state := "test-state-456"
	redirectURI := "http://localhost:8080/auth/callback/github"

	authURL := provider.GetAuthURL(state, redirectURI)

	// Parse the URL
	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	// Check base URL
	assert.Equal(t, "https", parsedURL.Scheme)
	assert.Equal(t, "github.com", parsedURL.Host)
	assert.Equal(t, "/login/oauth/authorize", parsedURL.Path)

	// Check query parameters
	query := parsedURL.Query()
	assert.Equal(t, "test-client-id", query.Get("client_id"))
	assert.Equal(t, redirectURI, query.Get("redirect_uri"))
	assert.Equal(t, state, query.Get("state"))
	assert.Contains(t, query.Get("scope"), "user:email")
	assert.Contains(t, query.Get("scope"), "read:user")
}

func TestMicrosoftOAuthProvider_GetProviderID(t *testing.T) {
	provider := NewMicrosoftOAuthProvider("test-client-id", "test-client-secret")
	assert.Equal(t, "microsoft", provider.GetProviderID())
}

func TestMicrosoftOAuthProvider_GetAuthURL(t *testing.T) {
	provider := NewMicrosoftOAuthProvider("test-client-id", "test-client-secret")

	state := "test-state-789"
	redirectURI := "http://localhost:8080/auth/callback/microsoft"

	authURL := provider.GetAuthURL(state, redirectURI)

	// Parse the URL
	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	// Check base URL
	assert.Equal(t, "https", parsedURL.Scheme)
	assert.Equal(t, "login.microsoftonline.com", parsedURL.Host)
	assert.Equal(t, "/common/oauth2/v2.0/authorize", parsedURL.Path)

	// Check query parameters
	query := parsedURL.Query()
	assert.Equal(t, "test-client-id", query.Get("client_id"))
	assert.Equal(t, redirectURI, query.Get("redirect_uri"))
	assert.Equal(t, "code", query.Get("response_type"))
	assert.Equal(t, state, query.Get("state"))
	assert.Contains(t, query.Get("scope"), "openid")
	assert.Contains(t, query.Get("scope"), "email")
	assert.Contains(t, query.Get("scope"), "profile")
}

// TestOAuthProviderExchangeCode tests error handling for ExchangeCode
// Since the actual implementation requires a real OAuth server, we test error cases
func TestOAuthProvider_ExchangeCodeErrors(t *testing.T) {
	tests := []struct {
		name        string
		provider    auth.OAuthProvider
		code        string
		redirectURI string
		wantErr     bool
	}{
		{
			name:        "google invalid code",
			provider:    NewGoogleOAuthProvider("test-client", "test-secret"),
			code:        "",
			redirectURI: "http://localhost",
			wantErr:     true,
		},
		{
			name:        "github invalid code",
			provider:    NewGitHubOAuthProvider("test-client", "test-secret"),
			code:        "",
			redirectURI: "http://localhost",
			wantErr:     true,
		},
		{
			name:        "microsoft invalid code",
			provider:    NewMicrosoftOAuthProvider("test-client", "test-secret"),
			code:        "",
			redirectURI: "http://localhost",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, err := tt.provider.ExchangeCode(ctx, tt.code, tt.redirectURI)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

// TestGitHubOAuthProvider_RefreshToken tests that GitHub doesn't support refresh tokens
func TestGitHubOAuthProvider_RefreshToken(t *testing.T) {
	provider := NewGitHubOAuthProvider("test-client", "test-secret")

	ctx := context.Background()
	_, err := provider.RefreshToken(ctx, "any-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not support refresh tokens")
}

// TestOAuthProviderURLValidation tests URL construction for all providers
func TestOAuthProviderURLValidation(t *testing.T) {
	tests := []struct {
		name        string
		provider    auth.OAuthProvider
		validateURL func(t *testing.T, parsedURL *url.URL)
	}{
		{
			name:     "google OAuth URLs",
			provider: NewGoogleOAuthProvider("client", "secret"),
			validateURL: func(t *testing.T, parsedURL *url.URL) {
				assert.Equal(t, "accounts.google.com", parsedURL.Host)
			},
		},
		{
			name:     "github OAuth URLs",
			provider: NewGitHubOAuthProvider("client", "secret"),
			validateURL: func(t *testing.T, parsedURL *url.URL) {
				assert.Equal(t, "github.com", parsedURL.Host)
			},
		},
		{
			name:     "microsoft OAuth URLs",
			provider: NewMicrosoftOAuthProvider("client", "secret"),
			validateURL: func(t *testing.T, parsedURL *url.URL) {
				assert.Equal(t, "login.microsoftonline.com", parsedURL.Host)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authURL := tt.provider.GetAuthURL("state", "http://localhost")
			parsedURL, err := url.Parse(authURL)
			require.NoError(t, err)

			// Validate HTTPS
			assert.Equal(t, "https", parsedURL.Scheme)

			// Provider-specific validation
			tt.validateURL(t, parsedURL)

			// Common query parameters
			query := parsedURL.Query()
			assert.NotEmpty(t, query.Get("client_id"))
			assert.NotEmpty(t, query.Get("redirect_uri"))
			assert.NotEmpty(t, query.Get("state"))
		})
	}
}

// TestExtractIDToken tests the ID token extraction helper
func TestExtractIDToken(t *testing.T) {
	tests := []struct {
		name     string
		token    *oauth2.Token
		expected string
	}{
		{
			name: "token with ID token",
			token: func() *oauth2.Token {
				tok := &oauth2.Token{
					AccessToken: "access",
				}
				return tok.WithExtra(map[string]interface{}{
					"id_token": "test-id-token",
				})
			}(),
			expected: "test-id-token",
		},
		{
			name: "token without ID token",
			token: &oauth2.Token{
				AccessToken: "access",
			},
			expected: "",
		},
		{
			name: "token with non-string ID token",
			token: func() *oauth2.Token {
				tok := &oauth2.Token{
					AccessToken: "access",
				}
				return tok.WithExtra(map[string]interface{}{
					"id_token": 12345,
				})
			}(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractIDToken(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}
