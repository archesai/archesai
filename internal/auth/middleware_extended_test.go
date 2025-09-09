package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestOptionalAuthMiddleware tests optional authentication middleware
func TestOptionalAuthMiddleware(t *testing.T) {
	t.Run("allows request without token", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()
		log := logger.NewTest()
		config := Config{
			JWTSecret:          "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		}
		service := NewService(mockRepo, mockUsersRepo, config, log)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := OptionalAuthMiddleware(service, log)
		handler := middleware(func(c echo.Context) error {
			// Should proceed without user context
			user := c.Get(string(AuthUserContextKey))
			assert.Nil(t, user)
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("sets context with valid token", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()
		log := logger.NewTest()
		config := Config{
			JWTSecret:          "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		}
		service := NewService(mockRepo, mockUsersRepo, config, log)

		e := echo.New()
		userID := uuid.New()
		sessionID := uuid.New()

		// Create a valid token
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID.String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			UserID:     userID,
			SessionID:  sessionID.String(),
			TokenType:  AccessTokenType,
			AuthMethod: AuthMethodPassword,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock the session validation
		mockRepo.EXPECT().GetSessionByToken(mock.Anything, mock.AnythingOfType("string")).Return(&Session{
			Id:        sessionID,
			UserId:    userID,
			Token:     "session-token",
			ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}, nil).Maybe()

		middleware := OptionalAuthMiddleware(service, log)
		handler := middleware(func(c echo.Context) error {
			user := c.Get(string(AuthUserContextKey))
			assert.NotNil(t, user)
			assert.Equal(t, userID, user)
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("continues with invalid token", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()
		log := logger.NewTest()
		config := Config{
			JWTSecret:          "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		}
		service := NewService(mockRepo, mockUsersRepo, config, log)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := OptionalAuthMiddleware(service, log)
		handler := middleware(func(c echo.Context) error {
			// Should proceed without user context
			user := c.Get(string(AuthUserContextKey))
			assert.Nil(t, user)
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// TestRequireRole tests role-based authorization middleware
func TestRequireRole(t *testing.T) {
	t.Run("allows access with required role", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Set up claims with admin role
		claims := &EnhancedClaims{
			UserID: uuid.New(),
			Roles:  []string{"admin", "user"},
		}
		c.Set(string(AuthClaimsContextKey), claims)

		middleware := RequireRole("admin")
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("denies access without required role", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Set up claims without admin role
		claims := &EnhancedClaims{
			UserID: uuid.New(),
			Roles:  []string{"user"},
		}
		c.Set(string(AuthClaimsContextKey), claims)

		middleware := RequireRole("admin")
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("returns unauthorized without claims", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := RequireRole("admin")
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})
}

// TestRequirePermission tests permission-based authorization middleware
func TestRequirePermission(t *testing.T) {
	t.Run("allows access with required permission", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &EnhancedClaims{
			UserID:      uuid.New(),
			Permissions: []string{"users:read", "users:write"},
		}
		c.Set(string(AuthClaimsContextKey), claims)

		middleware := RequirePermission("users:read")
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("denies access without required permission", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &EnhancedClaims{
			UserID:      uuid.New(),
			Permissions: []string{"users:read"},
		}
		c.Set(string(AuthClaimsContextKey), claims)

		middleware := RequirePermission("users:write")
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})
}

// TestRequireScope tests scope-based authorization middleware
func TestRequireScope(t *testing.T) {
	t.Run("allows access with required scope", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &EnhancedClaims{
			UserID: uuid.New(),
			Scopes: []string{"read", "write"},
		}
		c.Set(string(AuthClaimsContextKey), claims)

		middleware := RequireScope("read")
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("denies access without required scope", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &EnhancedClaims{
			UserID: uuid.New(),
			Scopes: []string{"read"},
		}
		c.Set(string(AuthClaimsContextKey), claims)

		middleware := RequireScope("write")
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})
}

// TestRequireOrganization tests organization-based authorization middleware
func TestRequireOrganization(t *testing.T) {
	t.Run("allows access with matching organization", func(t *testing.T) {
		e := echo.New()
		orgID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/orgs/"+orgID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("orgId")
		c.SetParamValues(orgID.String())

		claims := &EnhancedClaims{
			UserID:         uuid.New(),
			OrganizationID: orgID,
		}
		c.Set(string(AuthClaimsContextKey), claims)

		middleware := RequireOrganization()
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("denies access with different organization", func(t *testing.T) {
		e := echo.New()
		orgID := uuid.New()
		differentOrgID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/orgs/"+orgID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("orgId")
		c.SetParamValues(orgID.String())

		claims := &EnhancedClaims{
			UserID:         uuid.New(),
			OrganizationID: differentOrgID,
		}
		c.Set(string(AuthClaimsContextKey), claims)

		middleware := RequireOrganization()
		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		if err != nil {
			httpErr, ok := err.(*echo.HTTPError)
			if ok {
				assert.Equal(t, http.StatusForbidden, httpErr.Code)
			} else {
				assert.Error(t, err)
			}
		} else {
			// If no error, fail the test as we expect an error
			assert.Fail(t, "Expected an error but got none")
		}
	})
}

// TestGetUserFromContext tests getting user from context
func TestGetUserFromContext(t *testing.T) {
	t.Run("returns user when present", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		userID := uuid.New()
		c.Set(string(AuthUserContextKey), userID)

		result, ok := GetUserFromContext(c)
		assert.True(t, ok)
		assert.Equal(t, userID, result)
	})

	t.Run("returns zero UUID when not present", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result, ok := GetUserFromContext(c)
		assert.False(t, ok)
		assert.Equal(t, uuid.UUID{}, result)
	})
}

// TestGetClaimsFromContext tests getting claims from context
func TestGetClaimsFromContext(t *testing.T) {
	t.Run("returns claims when present", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &EnhancedClaims{
			UserID: uuid.New(),
		}
		c.Set(string(AuthClaimsContextKey), claims)

		result, ok := GetClaimsFromContext(c)
		assert.True(t, ok)
		assert.Equal(t, claims, result)
	})

	t.Run("returns nil when not present", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result, ok := GetClaimsFromContext(c)
		assert.False(t, ok)
		assert.Nil(t, result)
	})
}

// TestRateLimitMiddleware tests rate limiting middleware
func TestRateLimitMiddleware(t *testing.T) {
	t.Run("allows requests within limit", func(t *testing.T) {
		e := echo.New()
		middleware := RateLimitMiddleware(10, 1)

		// Make a few requests, should all succeed
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "127.0.0.1:1234"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := middleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "ok")
			})

			err := handler(c)
			assert.NoError(t, err)
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		e := echo.New()
		middleware := RateLimitMiddleware(2, 60)

		// First two requests should succeed
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "192.168.1.1:1234"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Response().Writer = rec

			handler := middleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "ok")
			})

			err := handler(c)
			assert.NoError(t, err)
		}

		// Third request should be rate limited
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusTooManyRequests, httpErr.Code)
	})
}

// TestSetRequestContextWithTimeout tests request context timeout middleware
func TestSetRequestContextWithTimeout(t *testing.T) {
	t.Run("sets context with timeout", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := SetRequestContextWithTimeout(100 * time.Millisecond)
		handler := middleware(func(c echo.Context) error {
			// Context should have a deadline
			ctx := c.Request().Context()
			deadline, ok := ctx.Deadline()
			assert.True(t, ok)
			assert.True(t, time.Until(deadline) <= 100*time.Millisecond)
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		// Create a context that we can cancel
		ctx, cancel := context.WithCancel(context.Background())
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := SetRequestContextWithTimeout(100 * time.Millisecond)
		handler := middleware(func(c echo.Context) error {
			// Cancel the context
			cancel()
			// Try to use the context
			select {
			case <-c.Request().Context().Done():
				return c.Request().Context().Err()
			default:
				return c.String(http.StatusOK, "ok")
			}
		})

		err := handler(c)
		assert.Error(t, err)
	})
}
