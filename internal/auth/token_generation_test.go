package auth

import (
	"context"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestService_GenerateTokensWithContext(t *testing.T) {
	// Setup
	repo := NewMockRepository(t)
	usersRepo := NewMockUsersRepository()
	config := Config{
		JWTSecret:          "test-secret-key-32-bytes-long!!!",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		SessionTokenExpiry: 30 * 24 * time.Hour,
	}
	service := NewService(repo, usersRepo, config, logger.NewTest())

	user := &users.User{
		Id:            uuid.New(),
		Email:         "test@example.com",
		Name:          "Test User",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	t.Run("generate basic tokens", func(t *testing.T) {
		tokens, err := service.generateTokens(user)
		if err != nil {
			t.Fatalf("generateTokens() error = %v", err)
		}

		if tokens.AccessToken == "" {
			t.Error("AccessToken is empty")
		}
		if tokens.RefreshToken == "" {
			t.Error("RefreshToken is empty")
		}
		if tokens.TokenType != "Bearer" {
			t.Errorf("TokenType = %v, want Bearer", tokens.TokenType)
		}
		if tokens.ExpiresIn != int64(config.AccessTokenExpiry.Seconds()) {
			t.Errorf("ExpiresIn = %v, want %v", tokens.ExpiresIn, int64(config.AccessTokenExpiry.Seconds()))
		}
	})

	t.Run("generate tokens with organization context", func(t *testing.T) {
		orgID := uuid.New()
		sessionID := uuid.New().String()
		ipAddress := "192.168.1.1"
		userAgent := "Mozilla/5.0"

		tokens, err := service.generateTokensWithContext(
			user,
			orgID,
			sessionID,
			ipAddress,
			userAgent,
			AuthMethodPassword,
			nil,
		)
		if err != nil {
			t.Fatalf("generateTokensWithContext() error = %v", err)
		}

		// Parse and verify access token
		token, err := jwt.ParseWithClaims(tokens.AccessToken, &EnhancedClaims{}, func(_ *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			t.Fatalf("Failed to parse access token: %v", err)
		}

		claims, ok := token.Claims.(*EnhancedClaims)
		if !ok {
			t.Fatal("Failed to cast claims to EnhancedClaims")
		}

		// Verify claims
		if claims.UserID != user.Id {
			t.Errorf("UserID = %v, want %v", claims.UserID, user.Id)
		}
		if claims.Email != string(user.Email) {
			t.Errorf("Email = %v, want %v", claims.Email, user.Email)
		}
		if claims.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", claims.OrganizationID, orgID)
		}
		if claims.SessionID != sessionID {
			t.Errorf("SessionID = %v, want %v", claims.SessionID, sessionID)
		}
		if claims.IPAddress != ipAddress {
			t.Errorf("IPAddress = %v, want %v", claims.IPAddress, ipAddress)
		}
		if claims.UserAgent != userAgent {
			t.Errorf("UserAgent = %v, want %v", claims.UserAgent, userAgent)
		}
		if claims.AuthMethod != AuthMethodPassword {
			t.Errorf("AuthMethod = %v, want %v", claims.AuthMethod, AuthMethodPassword)
		}
		if claims.TokenType != AccessTokenType {
			t.Errorf("TokenType = %v, want %v", claims.TokenType, AccessTokenType)
		}
	})

	t.Run("generate tokens with OAuth provider", func(t *testing.T) {
		provider := string(Google)
		tokens, err := service.generateTokensWithContext(
			user,
			uuid.Nil,
			"",
			"",
			"",
			AuthMethodOAuth,
			&provider,
		)
		if err != nil {
			t.Fatalf("generateTokensWithContext() error = %v", err)
		}

		// Parse and verify access token
		token, err := jwt.ParseWithClaims(tokens.AccessToken, &EnhancedClaims{}, func(_ *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			t.Fatalf("Failed to parse access token: %v", err)
		}

		claims, ok := token.Claims.(*EnhancedClaims)
		if !ok {
			t.Fatal("Failed to cast claims to EnhancedClaims")
		}

		if claims.Provider != provider {
			t.Errorf("Provider = %v, want %v", claims.Provider, provider)
		}
		if claims.AuthMethod != AuthMethodOAuth {
			t.Errorf("AuthMethod = %v, want %v", claims.AuthMethod, AuthMethodOAuth)
		}
	})

	t.Run("verify refresh token structure", func(t *testing.T) {
		sessionID := uuid.New().String()
		tokens, err := service.generateTokensWithContext(
			user,
			uuid.Nil,
			sessionID,
			"",
			"",
			AuthMethodPassword,
			nil,
		)
		if err != nil {
			t.Fatalf("generateTokensWithContext() error = %v", err)
		}

		// Parse refresh token
		token, err := jwt.ParseWithClaims(tokens.RefreshToken, &RefreshClaims{}, func(_ *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			t.Fatalf("Failed to parse refresh token: %v", err)
		}

		claims, ok := token.Claims.(*RefreshClaims)
		if !ok {
			t.Fatal("Failed to cast claims to RefreshClaims")
		}

		if claims.UserID != user.Id {
			t.Errorf("UserID = %v, want %v", claims.UserID, user.Id)
		}
		if claims.TokenType != RefreshTokenType {
			t.Errorf("TokenType = %v, want %v", claims.TokenType, RefreshTokenType)
		}
		if claims.SessionID != sessionID {
			t.Errorf("SessionID = %v, want %v", claims.SessionID, sessionID)
		}
		if claims.AuthMethod != AuthMethodPassword {
			t.Errorf("AuthMethod = %v, want %v", claims.AuthMethod, AuthMethodPassword)
		}
	})
}

func TestService_ValidateEnhancedToken(t *testing.T) {
	// Setup
	repo := NewMockRepository(t)
	usersRepo := NewMockUsersRepository()
	config := Config{
		JWTSecret:          "test-secret-key-32-bytes-long!!!",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	service := NewService(repo, usersRepo, config, logger.NewTest())

	user := &users.User{
		Id:            uuid.New(),
		Email:         "test@example.com",
		Name:          "Test User",
		EmailVerified: true,
	}

	t.Run("validate valid token", func(t *testing.T) {
		// Generate token
		tokens, err := service.generateTokens(user)
		if err != nil {
			t.Fatalf("generateTokens() error = %v", err)
		}

		// Validate token
		claims, err := service.ValidateToken(tokens.AccessToken)
		if err != nil {
			t.Fatalf("ValidateToken() error = %v", err)
		}

		if claims.UserID != user.Id {
			t.Errorf("UserID = %v, want %v", claims.UserID, user.Id)
		}
		if claims.Email != string(user.Email) {
			t.Errorf("Email = %v, want %v", claims.Email, user.Email)
		}
	})

	t.Run("validate expired token", func(t *testing.T) {
		// Create expired token
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Subject:   user.Id.String(),
			},
			UserID:    user.Id,
			Email:     string(user.Email),
			TokenType: AccessTokenType,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Validate should fail
		_, err = service.ValidateToken(tokenString)
		if err != ErrTokenExpired {
			t.Errorf("ValidateToken() error = %v, want %v", err, ErrTokenExpired)
		}
	})

	t.Run("validate invalid signature", func(t *testing.T) {
		// Create token with different secret
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   user.Id.String(),
			},
			UserID:    user.Id,
			Email:     string(user.Email),
			TokenType: AccessTokenType,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("wrong-secret-key"))
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Validate should fail
		_, err = service.ValidateToken(tokenString)
		if err != ErrInvalidToken {
			t.Errorf("ValidateToken() error = %v, want %v", err, ErrInvalidToken)
		}
	})

	t.Run("validate malformed token", func(t *testing.T) {
		_, err := service.ValidateToken("not.a.token")
		if err != ErrInvalidToken {
			t.Errorf("ValidateToken() error = %v, want %v", err, ErrInvalidToken)
		}
	})
}

func TestService_RefreshTokenWithEnhancedClaims(t *testing.T) {
	// Setup
	repo := NewMockRepository(t)
	usersRepo := NewMockUsersRepository()
	config := Config{
		JWTSecret:          "test-secret-key-32-bytes-long!!!",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		SessionTokenExpiry: 30 * 24 * time.Hour,
	}
	service := NewService(repo, usersRepo, config, logger.NewTest())

	user := &users.User{
		Id:            uuid.New(),
		Email:         "test@example.com",
		Name:          "Test User",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Add user to mock repository
	usersRepo.users[user.Id] = user

	t.Run("refresh with valid refresh token", func(t *testing.T) {
		ctx := context.Background()
		sessionID := uuid.New()

		// Create a session
		session := &Session{
			Id:                   sessionID,
			UserId:               user.Id,
			Token:                "refresh-token",
			ExpiresAt:            time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			ActiveOrganizationId: uuid.New(),
			IpAddress:            "192.168.1.1",
			UserAgent:            "Test Agent",
		}

		// Mock GetSession expectation for potential session lookup
		repo.EXPECT().GetSession(mock.Anything, sessionID).Return(session, nil).Maybe()

		// Generate initial tokens with session
		tokens, err := service.generateTokensWithContext(
			user,
			session.ActiveOrganizationId,
			sessionID.String(),
			session.IpAddress,
			session.UserAgent,
			AuthMethodPassword,
			nil,
		)
		if err != nil {
			t.Fatalf("generateTokensWithContext() error = %v", err)
		}

		// Refresh token
		newTokens, err := service.RefreshToken(ctx, tokens.RefreshToken)
		if err != nil {
			t.Fatalf("RefreshToken() error = %v", err)
		}

		if newTokens.AccessToken == "" {
			t.Error("New AccessToken is empty")
		}
		if newTokens.RefreshToken == "" {
			t.Error("New RefreshToken is empty")
		}
		if newTokens.AccessToken == tokens.AccessToken {
			t.Error("New AccessToken should be different from old one")
		}
		if newTokens.RefreshToken == tokens.RefreshToken {
			t.Error("New RefreshToken should be different from old one")
		}
	})

	t.Run("refresh with expired refresh token", func(t *testing.T) {
		ctx := context.Background()

		// Create expired refresh token
		refreshClaims := &RefreshClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
				Subject:   user.Id.String(),
			},
			UserID:     user.Id,
			TokenType:  RefreshTokenType,
			SessionID:  uuid.New().String(),
			AuthMethod: AuthMethodPassword,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Refresh should fail
		_, err = service.RefreshToken(ctx, tokenString)
		if err != ErrTokenExpired {
			t.Errorf("RefreshToken() error = %v, want %v", err, ErrTokenExpired)
		}
	})

	t.Run("refresh with non-refresh token", func(t *testing.T) {
		ctx := context.Background()

		// Create access token (not refresh token)
		accessClaims := &RefreshClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   user.Id.String(),
			},
			UserID:     user.Id,
			TokenType:  AccessTokenType, // Wrong token type
			SessionID:  uuid.New().String(),
			AuthMethod: AuthMethodPassword,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Refresh should fail
		_, err = service.RefreshToken(ctx, tokenString)
		if err != ErrInvalidToken {
			t.Errorf("RefreshToken() error = %v, want %v", err, ErrInvalidToken)
		}
	})

	t.Run("refresh with non-existent user", func(t *testing.T) {
		ctx := context.Background()
		nonExistentUserID := uuid.New()

		// Create refresh token for non-existent user
		refreshClaims := &RefreshClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   nonExistentUserID.String(),
			},
			UserID:     nonExistentUserID,
			TokenType:  RefreshTokenType,
			SessionID:  uuid.New().String(),
			AuthMethod: AuthMethodPassword,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Refresh should fail
		_, err = service.RefreshToken(ctx, tokenString)
		if err != ErrUserNotFound {
			t.Errorf("RefreshToken() error = %v, want %v", err, ErrUserNotFound)
		}
	})
}

func TestService_ValidateLegacyToken(t *testing.T) {
	// Setup
	repo := NewMockRepository(t)
	usersRepo := NewMockUsersRepository()
	config := Config{
		JWTSecret:         "test-secret-key-32-bytes-long!!!",
		AccessTokenExpiry: 15 * time.Minute,
	}
	service := NewService(repo, usersRepo, config, logger.NewTest())

	userID := uuid.New()
	email := "legacy@example.com"

	t.Run("validate legacy token", func(t *testing.T) {
		// Create legacy token
		claims := &Claims{
			UserID: userID,
			Email:  email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Validate legacy token
		validatedClaims, err := service.ValidateLegacyToken(tokenString)
		if err != nil {
			t.Fatalf("ValidateLegacyToken() error = %v", err)
		}

		if validatedClaims.UserID != userID {
			t.Errorf("UserID = %v, want %v", validatedClaims.UserID, userID)
		}
		if validatedClaims.Email != email {
			t.Errorf("Email = %v, want %v", validatedClaims.Email, email)
		}
	})
}

func BenchmarkGenerateTokens(b *testing.B) {
	// Setup
	repo := NewMockRepository(b)
	usersRepo := NewMockUsersRepository()
	config := Config{
		JWTSecret:          "test-secret-key-32-bytes-long!!!",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	service := NewService(repo, usersRepo, config, logger.NewTest())

	user := &users.User{
		Id:            uuid.New(),
		Email:         "bench@example.com",
		Name:          "Bench User",
		EmailVerified: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.generateTokens(user)
	}
}

func BenchmarkValidateToken(b *testing.B) {
	// Setup
	repo := NewMockRepository(b)
	usersRepo := NewMockUsersRepository()
	config := Config{
		JWTSecret:         "test-secret-key-32-bytes-long!!!",
		AccessTokenExpiry: 15 * time.Minute,
	}
	service := NewService(repo, usersRepo, config, logger.NewTest())

	user := &users.User{
		Id:            uuid.New(),
		Email:         "bench@example.com",
		Name:          "Bench User",
		EmailVerified: true,
	}

	tokens, _ := service.generateTokens(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateToken(tokens.AccessToken)
	}
}
