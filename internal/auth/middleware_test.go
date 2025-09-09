package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		// Create mocks
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()

		// Create a test user
		userID := uuid.New()
		testUser := &users.User{
			Id:    userID,
			Email: "test@example.com",
			Name:  "Test User",
		}
		mockUsersRepo.users[userID] = testUser

		// No need to setup session expectations - middleware doesn't check sessions

		service := &Service{
			repo:      mockRepo,
			usersRepo: mockUsersRepo,
			jwtSecret: []byte("test-secret"),
			logger:    logger.NewTest(),
			config: Config{
				JWTSecret:          "test-secret",
				AccessTokenExpiry:  15 * time.Minute,
				RefreshTokenExpiry: 7 * 24 * time.Hour,
			},
		}

		middleware := Middleware(service, logger.NewTest())

		// Create a valid JWT token with proper EnhancedClaims structure
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
			UserID:     userID,
			Email:      "test@example.com",
			TokenType:  AccessTokenType,
			AuthMethod: AuthMethodPassword,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret"))

		// Create Echo context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that will be called if middleware passes
		handler := func(c echo.Context) error {
			// Check that user context was set
			user := c.Get(string(AuthUserContextKey))
			if user == nil {
				t.Error("Expected user in context")
			}
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("missing token", func(t *testing.T) {
		// Create mocks
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()

		service := &Service{
			repo:      mockRepo,
			usersRepo: mockUsersRepo,
			jwtSecret: []byte("test-secret"),
			logger:    logger.NewTest(),
			config: Config{
				JWTSecret: "test-secret",
			},
		}

		middleware := Middleware(service, logger.NewTest())

		// Create Echo context without auth header
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that should not be called
		handler := func(c echo.Context) error {
			t.Error("Handler should not be called")
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert - should return unauthorized
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		// Create mocks
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()

		service := &Service{
			repo:      mockRepo,
			usersRepo: mockUsersRepo,
			jwtSecret: []byte("test-secret"),
			logger:    logger.NewTest(),
			config: Config{
				JWTSecret: "test-secret",
			},
		}

		middleware := Middleware(service, logger.NewTest())

		// Create Echo context with invalid token
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that should not be called
		handler := func(c echo.Context) error {
			t.Error("Handler should not be called")
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert - should return unauthorized
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	t.Run("expired token", func(t *testing.T) {
		// Create mocks
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()

		service := &Service{
			repo:      mockRepo,
			usersRepo: mockUsersRepo,
			jwtSecret: []byte("test-secret"),
			logger:    logger.NewTest(),
			config: Config{
				JWTSecret: "test-secret",
			},
		}

		middleware := Middleware(service, logger.NewTest())

		// Create an expired JWT token
		userID := uuid.New()
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Subject:   userID.String(),
			},
			UserID:     userID,
			Email:      "test@example.com",
			TokenType:  AccessTokenType,
			AuthMethod: AuthMethodPassword,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret"))

		// Create Echo context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that should not be called
		handler := func(c echo.Context) error {
			t.Error("Handler should not be called")
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert - should return unauthorized
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	t.Run("session not found", func(t *testing.T) {
		// Create mocks
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()

		// No session expectations needed - middleware doesn't check sessions

		service := &Service{
			repo:      mockRepo,
			usersRepo: mockUsersRepo,
			jwtSecret: []byte("test-secret"),
			logger:    logger.NewTest(),
			config: Config{
				JWTSecret: "test-secret",
			},
		}

		middleware := Middleware(service, logger.NewTest())

		// Create a valid JWT token but session doesn't exist
		userID := uuid.New()
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
			UserID:     userID,
			Email:      "test@example.com",
			TokenType:  AccessTokenType,
			AuthMethod: AuthMethodPassword,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret"))

		// Create Echo context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that should not be called
		handler := func(c echo.Context) error {
			t.Error("Handler should not be called")
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert - should return unauthorized
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})
}

func TestRequireAuthMiddleware(t *testing.T) {
	t.Run("with user context", func(t *testing.T) {
		// Create mocks
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()

		// Create a test user
		userID := uuid.New()
		testUser := &users.User{
			Id:    userID,
			Email: "test@example.com",
			Name:  "Test User",
		}
		mockUsersRepo.users[userID] = testUser

		service := &Service{
			repo:      mockRepo,
			usersRepo: mockUsersRepo,
			jwtSecret: []byte("test-secret"),
			logger:    logger.NewTest(),
			config: Config{
				JWTSecret: "test-secret",
			},
		}

		middleware := Middleware(service, logger.NewTest())

		// Create a valid JWT token
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
			UserID:     userID,
			Email:      "test@example.com",
			TokenType:  AccessTokenType,
			AuthMethod: AuthMethodPassword,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test-secret"))

		// Create Echo context with token
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that should be called
		handler := func(c echo.Context) error {
			user := c.Get(string(AuthUserContextKey))
			if user == nil {
				t.Error("Expected user in context")
			}
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("without user context", func(t *testing.T) {
		// Create mocks
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()

		service := &Service{
			repo:      mockRepo,
			usersRepo: mockUsersRepo,
			jwtSecret: []byte("test-secret"),
			logger:    logger.NewTest(),
			config: Config{
				JWTSecret: "test-secret",
			},
		}

		middleware := Middleware(service, logger.NewTest())

		// Create Echo context without user
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that should not be called
		handler := func(c echo.Context) error {
			t.Error("Handler should not be called")
			return c.String(http.StatusOK, "success")
		}

		// Execute middleware
		err := middleware(handler)(c)

		// Assert - should return unauthorized
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})
}
