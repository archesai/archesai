package oauth

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/archesai/archesai/internal/logger"
)

func TestMicrosoftOAuthProvider_GetProviderID(t *testing.T) {
	provider := NewMicrosoftProvider(
		"test-client-id",
		"test-client-secret",
		"http://localhost:8080/auth/callback/microsoft",
		logger.NewTest(),
	)
	assert.Equal(t, "microsoft", provider.GetProviderID())
}

func TestMicrosoftOAuthProvider_GetAuthURL(t *testing.T) {
	provider := NewMicrosoftProvider(
		"test-client-id",
		"test-client-secret",
		"http://localhost:8080/auth/callback/microsoft",
		logger.NewTest(),
	)

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

func TestMicrosoftOAuthProvider_ExchangeCodeError(t *testing.T) {
	provider := NewMicrosoftProvider(
		"test-client",
		"test-secret",
		"http://localhost:8080/auth/callback/microsoft",
		logger.NewTest(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := provider.ExchangeCode(ctx, "", "http://localhost")
	assert.Error(t, err)
}

func TestMicrosoftOAuthProvider_URLValidation(t *testing.T) {
	provider := NewMicrosoftProvider(
		"client",
		"secret",
		"http://localhost:8080/auth/callback/microsoft",
		logger.NewTest(),
	)

	authURL := provider.GetAuthURL("state", "http://localhost")
	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	// Validate HTTPS
	assert.Equal(t, "https", parsedURL.Scheme)
	assert.Equal(t, "login.microsoftonline.com", parsedURL.Host)

	// Common query parameters
	query := parsedURL.Query()
	assert.NotEmpty(t, query.Get("client_id"))
	assert.NotEmpty(t, query.Get("redirect_uri"))
	assert.NotEmpty(t, query.Get("state"))
}
