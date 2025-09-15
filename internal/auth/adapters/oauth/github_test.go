package oauth

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestGitHubOAuthProvider_ExchangeCodeError(t *testing.T) {
	provider := NewGitHubOAuthProvider("test-client", "test-secret")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := provider.ExchangeCode(ctx, "", "http://localhost")
	assert.Error(t, err)
}

// TestGitHubOAuthProvider_RefreshToken tests that GitHub doesn't support refresh tokens
func TestGitHubOAuthProvider_RefreshToken(t *testing.T) {
	provider := NewGitHubOAuthProvider("test-client", "test-secret")

	ctx := context.Background()
	_, err := provider.RefreshToken(ctx, "any-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not support refresh tokens")
}

func TestGitHubOAuthProvider_URLValidation(t *testing.T) {
	provider := NewGitHubOAuthProvider("client", "secret")

	authURL := provider.GetAuthURL("state", "http://localhost")
	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	// Validate HTTPS
	assert.Equal(t, "https", parsedURL.Scheme)
	assert.Equal(t, "github.com", parsedURL.Host)

	// Common query parameters
	query := parsedURL.Query()
	assert.NotEmpty(t, query.Get("client_id"))
	assert.NotEmpty(t, query.Get("redirect_uri"))
	assert.NotEmpty(t, query.Get("state"))
}
