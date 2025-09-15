package auth

import (
	"context"
	"net/url"
	"testing"
	"time"

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

func TestGoogleOAuthProvider_ExchangeCodeError(t *testing.T) {
	provider := NewGoogleOAuthProvider("test-client", "test-secret")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := provider.ExchangeCode(ctx, "", "http://localhost")
	assert.Error(t, err)
}

func TestGoogleOAuthProvider_URLValidation(t *testing.T) {
	provider := NewGoogleOAuthProvider("client", "secret")

	authURL := provider.GetAuthURL("state", "http://localhost")
	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	// Validate HTTPS
	assert.Equal(t, "https", parsedURL.Scheme)
	assert.Equal(t, "accounts.google.com", parsedURL.Host)

	// Common query parameters
	query := parsedURL.Query()
	assert.NotEmpty(t, query.Get("client_id"))
	assert.NotEmpty(t, query.Get("redirect_uri"))
	assert.NotEmpty(t, query.Get("state"))
}

func TestGoogleOAuthProvider_ExtractIDToken(t *testing.T) {
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
