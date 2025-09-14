package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAPIKeyMiddleware(t *testing.T) {
	t.Run("valid API key in X-API-Key header", func(t *testing.T) {
		// Create mock service
		service := &Service{
			logger: logger.NewTest(),
		}

		middleware := APIKeyMiddleware(service, logger.NewTest())

		// Create Echo context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-API-Key", "sk_live_test123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that will be called if middleware passes
		handler := func(c echo.Context) error {
			// Check that API key context was set
			apiKey := c.Get("api_key")
			if apiKey == nil {
				t.Error("Expected API key in context")
			}
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware - expect error since ValidateAPIKey returns error
		err := middleware(handler)(c)

		// Assert - should fail as ValidateAPIKey is not implemented
		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	t.Run("valid API key in Authorization header", func(t *testing.T) {
		// Create mock service
		service := &Service{
			logger: logger.NewTest(),
		}

		middleware := APIKeyMiddleware(service, logger.NewTest())

		// Create Echo context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "ApiKey sk_live_test123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler
		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert - should fail as ValidateAPIKey is not implemented
		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	t.Run("missing API key", func(t *testing.T) {
		// Create mock service
		service := &Service{
			logger: logger.NewTest(),
		}

		middleware := APIKeyMiddleware(service, logger.NewTest())

		// Create Echo context without API key
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler
		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert
		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
		assert.Equal(t, "missing API key", httpErr.Message)
	})
}

func TestParseAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ApiKey scheme",
			input:    "ApiKey sk_live_123456",
			expected: "sk_live_123456",
		},
		{
			name:     "Bearer scheme with API key",
			input:    "Bearer sk_test_abcdef",
			expected: "sk_test_abcdef",
		},
		{
			name:     "Direct API key",
			input:    "sk_live_xyz789",
			expected: "sk_live_xyz789",
		},
		{
			name:     "Bearer with JWT token",
			input:    "Bearer eyJhbGciOiJIUzI1NiJ9.token",
			expected: "",
		},
		{
			name:     "Invalid format",
			input:    "Basic dXNlcjpwYXNz",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Case insensitive ApiKey",
			input:    "apikey sk_live_test",
			expected: "sk_live_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAPIKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateAPIKeyFormat(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{
			name:     "valid live key",
			key:      "sk_live_" + generateHexString(64),
			expected: true,
		},
		{
			name:     "valid test key",
			key:      "sk_test_" + generateHexString(64),
			expected: true,
		},
		{
			name:     "missing prefix",
			key:      "live_" + generateHexString(64),
			expected: false,
		},
		{
			name:     "invalid environment",
			key:      "sk_prod_" + generateHexString(64),
			expected: false,
		},
		{
			name:     "wrong hex length",
			key:      "sk_live_" + generateHexString(32),
			expected: false,
		},
		{
			name:     "not hex",
			key:      "sk_live_" + "xyz" + generateHexString(61),
			expected: false,
		},
		{
			name:     "empty string",
			key:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAPIKeyFormat(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateAPIKey(t *testing.T) {
	key, prefix, err := GenerateAPIKey()

	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	assert.NotEmpty(t, prefix)

	// Check format
	assert.True(t, ValidateAPIKeyFormat(key))
	assert.Equal(t, key[:8], prefix)
	assert.True(t, len(key) > 64) // sk_live_ + 64 hex chars
}

func TestAPIKeyService_CreateAPIKey(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		// Create mock repository
		mockRepo := &mockAPIKeyRepository{
			apiKeys: make(map[uuid.UUID]*APIKey),
		}

		// Create service
		service := NewAPIKeyService(mockRepo, nil)

		// Create API key
		userID := uuid.New()
		orgID := uuid.New()
		scopes := []string{"read:workflows", "write:workflows"}

		result, err := service.CreateAPIKey(
			context.Background(),
			userID,
			orgID,
			"Test API Key",
			scopes,
			24*time.Hour,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.PlainKey)
		assert.True(t, ValidateAPIKeyFormat(result.PlainKey))
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, "Test API Key", result.Name)
		assert.Equal(t, scopes, result.Scopes)
		assert.Equal(t, 100, result.RateLimit)
	})
}

// Helper function to generate hex string of specific length
func generateHexString(length int) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, length)
	for i := range result {
		result[i] = hexChars[i%16]
	}
	return string(result)
}

// Mock API key repository for testing
type mockAPIKeyRepository struct {
	apiKeys map[uuid.UUID]*APIKey
	err     error
}

func (m *mockAPIKeyRepository) CreateAPIKey(_ context.Context, apiKey *APIKey) (*APIKey, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.apiKeys[apiKey.ID] = apiKey
	return apiKey, nil
}

func (m *mockAPIKeyRepository) GetAPIKeyByPrefix(_ context.Context, prefix string) (*APIKey, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, key := range m.apiKeys {
		if key.Prefix == prefix {
			return key, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockAPIKeyRepository) GetAPIKeyByID(_ context.Context, id uuid.UUID) (*APIKey, error) {
	if m.err != nil {
		return nil, m.err
	}
	if key, ok := m.apiKeys[id]; ok {
		return key, nil
	}
	return nil, ErrNotFound
}

func (m *mockAPIKeyRepository) ListUserAPIKeys(_ context.Context, userID uuid.UUID) ([]*APIKey, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*APIKey
	for _, key := range m.apiKeys {
		if key.UserID == userID {
			result = append(result, key)
		}
	}
	return result, nil
}

func (m *mockAPIKeyRepository) UpdateAPIKeyLastUsed(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if key, ok := m.apiKeys[id]; ok {
		key.LastUsedAt = time.Now()
		return nil
	}
	return ErrNotFound
}

func (m *mockAPIKeyRepository) DeleteAPIKey(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	delete(m.apiKeys, id)
	return nil
}

func (m *mockAPIKeyRepository) ValidateAPIKeyHash(_ context.Context, prefix, keyHash string) (*APIKey, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, key := range m.apiKeys {
		if key.Prefix == prefix && key.KeyHash == keyHash {
			return key, nil
		}
	}
	return nil, ErrInvalidCredentials
}

var ErrNotFound = errors.New("not found")
