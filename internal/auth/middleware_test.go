package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a test service with mock repository
	mockRepo := &MockRepository{
		users:    make(map[uuid.UUID]*User),
		sessions: make(map[uuid.UUID]*Session),
		accounts: make(map[uuid.UUID]*Account),
	}

	service := &Service{
		repo:      mockRepo,
		jwtSecret: []byte("test-secret"),
		config: Config{
			JWTSecret:          "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
	}

	middleware := Middleware(service, logger.NewTest())

	t.Run("valid token", func(t *testing.T) {
		// Create a test user
		userID := uuid.New()
		user := &User{
			Id:            userID,
			Email:         "test@example.com",
			Name:          "Test User",
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		mockRepo.users[userID] = user

		// Generate a valid token
		claims := &Claims{
			UserID: userID,
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(service.jwtSecret)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		// Create test request with token
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that checks if user is in context
		handler := func(c echo.Context) error {
			userID := c.Get(string(AuthUserContextKey))
			if userID == nil {
				t.Error("Expected user ID in context, got nil")
			}
			claims := c.Get(string(AuthClaimsContextKey))
			if claims == nil {
				t.Error("Expected claims in context, got nil")
			}
			return c.String(http.StatusOK, "OK")
		}

		// Apply middleware and call handler
		h := middleware(handler)
		err = h(c)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}

		h := middleware(handler)
		err := h(c)

		if err == nil {
			t.Error("Expected error for missing token")
		}
		httpErr, ok := err.(*echo.HTTPError)
		if !ok {
			t.Errorf("Expected echo.HTTPError, got %T", err)
		} else if httpErr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", httpErr.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}

		h := middleware(handler)
		err := h(c)

		if err == nil {
			t.Error("Expected error for invalid token")
		}
		httpErr, ok := err.(*echo.HTTPError)
		if !ok {
			t.Errorf("Expected echo.HTTPError, got %T", err)
		} else if httpErr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", httpErr.Code)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		// Create an expired token
		claims := &Claims{
			UserID: uuid.New(),
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Subject:   uuid.New().String(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(service.jwtSecret)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}

		h := middleware(handler)
		err = h(c)

		if err == nil {
			t.Error("Expected error for expired token")
		}
		httpErr, ok := err.(*echo.HTTPError)
		if !ok {
			t.Errorf("Expected echo.HTTPError, got %T", err)
		} else if httpErr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", httpErr.Code)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		// Create a token for non-existent user
		userID := uuid.New()
		claims := &Claims{
			UserID: userID,
			Email:  "nonexistent@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(service.jwtSecret)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}

		h := middleware(handler)
		err = h(c)

		if err == nil {
			t.Error("Expected error for non-existent user")
		}
		httpErr, ok := err.(*echo.HTTPError)
		if !ok {
			t.Errorf("Expected echo.HTTPError, got %T", err)
		} else if httpErr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", httpErr.Code)
		}
	})
}

func TestOptionalAuthMiddleware(t *testing.T) {
	// Create a test service with mock repository
	mockRepo := &MockRepository{
		users:    make(map[uuid.UUID]*User),
		sessions: make(map[uuid.UUID]*Session),
		accounts: make(map[uuid.UUID]*Account),
	}

	service := &Service{
		repo:      mockRepo,
		jwtSecret: []byte("test-secret"),
		config: Config{
			JWTSecret:          "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
	}

	middleware := OptionalAuthMiddleware(service, logger.NewTest())

	t.Run("with valid token", func(t *testing.T) {
		// Create a test user
		userID := uuid.New()
		user := &User{
			Id:            userID,
			Email:         "test@example.com",
			Name:          "Test User",
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		mockRepo.users[userID] = user

		// Generate a valid token
		claims := &Claims{
			UserID: userID,
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(service.jwtSecret)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		// Create test request with token
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that checks if user is in context
		handler := func(c echo.Context) error {
			user := c.Get(string(UserContextKey))
			if user == nil {
				t.Error("Expected user in context, got nil")
			}
			return c.String(http.StatusOK, "OK")
		}

		// Apply middleware and call handler
		h := middleware(handler)
		err = h(c)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("without token", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that checks user is NOT in context
		handler := func(c echo.Context) error {
			user := c.Get(string(UserContextKey))
			if user != nil {
				t.Error("Expected no user in context, got user")
			}
			return c.String(http.StatusOK, "OK")
		}

		// Apply middleware and call handler
		h := middleware(handler)
		err := h(c)

		// Should not error - optional auth allows missing token
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("with invalid token", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a handler that checks user is NOT in context
		handler := func(c echo.Context) error {
			user := c.Get(string(UserContextKey))
			if user != nil {
				t.Error("Expected no user in context with invalid token")
			}
			return c.String(http.StatusOK, "OK")
		}

		// Apply middleware and call handler
		h := middleware(handler)
		err := h(c)

		// Should not error - optional auth allows invalid token
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})
}

func TestGetUserFromContext(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() echo.Context
		wantUser bool
	}{
		{
			name: "user in context",
			setupCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set(string(AuthUserContextKey), "user-id-123")
				return c
			},
			wantUser: true,
		},
		{
			name: "no user in context",
			setupCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return c
			},
			wantUser: false,
		},
		{
			name: "wrong type in context",
			setupCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set(string(AuthUserContextKey), 123) // wrong type
				return c
			},
			wantUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			userID, ok := GetUserFromContext(c)

			if tt.wantUser {
				if !ok {
					t.Error("Expected user ID, got false")
				}
				if userID == "" {
					t.Error("Expected non-empty user ID")
				}
			} else {
				if ok {
					t.Error("Expected no user ID, got true")
				}
			}
		})
	}
}

func TestGetClaimsFromContext(t *testing.T) {
	tests := []struct {
		name       string
		setupCtx   func() echo.Context
		wantClaims bool
	}{
		{
			name: "claims in context",
			setupCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				claims := &Claims{
					UserID: uuid.New(),
					Email:  "test@example.com",
				}
				c.Set(string(AuthClaimsContextKey), claims)
				return c
			},
			wantClaims: true,
		},
		{
			name: "no claims in context",
			setupCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return c
			},
			wantClaims: false,
		},
		{
			name: "wrong type in context",
			setupCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set(string(AuthClaimsContextKey), "not claims")
				return c
			},
			wantClaims: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			claims, ok := GetClaimsFromContext(c)

			if tt.wantClaims {
				if !ok {
					t.Error("Expected claims, got false")
				}
				if claims == nil {
					t.Error("Expected non-nil claims")
				}
			} else {
				if ok {
					t.Error("Expected no claims, got true")
				}
			}
		})
	}
}
