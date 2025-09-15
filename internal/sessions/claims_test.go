package sessions

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testEmail = "test@example.com"

func TestEnhancedClaims_Creation(t *testing.T) {
	userID := uuid.New()
	email := testEmail
	now := time.Now()

	t.Run("basic claims creation", func(t *testing.T) {
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "archesai",
				Subject:   userID.String(),
				ID:        uuid.New().String(),
			},
			UserID:       userID,
			Email:        email,
			TokenType:    AccessTokenType,
			AuthMethod:   AuthMethodPassword,
			Features:     make(map[string]bool),
			CustomClaims: make(map[string]interface{}),
		}

		if claims.UserID != userID {
			t.Errorf("UserID = %v, want %v", claims.UserID, userID)
		}
		if claims.Email != email {
			t.Errorf("Email = %v, want %v", claims.Email, email)
		}
		if claims.TokenType != AccessTokenType {
			t.Errorf("TokenType = %v, want %v", claims.TokenType, AccessTokenType)
		}
		if claims.AuthMethod != AuthMethodPassword {
			t.Errorf("AuthMethod = %v, want %v", claims.AuthMethod, AuthMethodPassword)
		}
	})

	t.Run("with expiry", func(t *testing.T) {
		duration := 15 * time.Minute
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Subject:   userID.String(),
			},
			UserID: userID,
			Email:  email,
		}

		expectedExpiry := now.Add(duration)
		actualExpiry := claims.ExpiresAt.Time

		// Allow 1 second tolerance
		if actualExpiry.Sub(expectedExpiry) > time.Second {
			t.Errorf("ExpiresAt off by %v", actualExpiry.Sub(expectedExpiry))
		}
	})

	t.Run("with user info", func(t *testing.T) {
		name := "Test User"
		avatarURL := "https://example.com/avatar.jpg"
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID.String(),
			},
			UserID:        userID,
			Email:         email,
			Name:          name,
			AvatarURL:     avatarURL,
			EmailVerified: true,
		}

		if claims.Name != name {
			t.Errorf("Name = %v, want %v", claims.Name, name)
		}
		if claims.AvatarURL != avatarURL {
			t.Errorf("AvatarURL = %v, want %v", claims.AvatarURL, avatarURL)
		}
		if !claims.EmailVerified {
			t.Error("EmailVerified should be true")
		}
	})

	t.Run("with organization", func(t *testing.T) {
		orgID := uuid.New()
		orgName := "Test Org"
		role := "admin"

		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID.String(),
			},
			UserID:           userID,
			Email:            email,
			OrganizationID:   orgID,
			OrganizationName: orgName,
			OrganizationRole: role,
		}

		if claims.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", claims.OrganizationID, orgID)
		}
		if claims.OrganizationName != orgName {
			t.Errorf("OrganizationName = %v, want %v", claims.OrganizationName, orgName)
		}
		if claims.OrganizationRole != role {
			t.Errorf("OrganizationRole = %v, want %v", claims.OrganizationRole, role)
		}
	})

	t.Run("with multiple organizations", func(t *testing.T) {
		orgs := []OrganizationClaim{
			{
				ID:          uuid.New(),
				Name:        "Org 1",
				Role:        "admin",
				Permissions: []string{"read", "write"},
			},
			{
				ID:   uuid.New(),
				Name: "Org 2",
				Role: "member",
			},
		}

		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID.String(),
			},
			UserID:        userID,
			Email:         email,
			Organizations: orgs,
		}

		if len(claims.Organizations) != 2 {
			t.Errorf("Organizations length = %v, want 2", len(claims.Organizations))
		}
		if claims.Organizations[0].Name != "Org 1" {
			t.Errorf("First org name = %v, want Org 1", claims.Organizations[0].Name)
		}
	})

	t.Run("with roles and permissions", func(t *testing.T) {
		roles := []string{"admin", "developer"}
		permissions := []string{"read", "write", "delete"}

		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID.String(),
			},
			UserID:      userID,
			Email:       email,
			Roles:       roles,
			Permissions: permissions,
		}

		if len(claims.Roles) != 2 {
			t.Errorf("Roles length = %v, want 2", len(claims.Roles))
		}
		if len(claims.Permissions) != 3 {
			t.Errorf("Permissions length = %v, want 3", len(claims.Permissions))
		}
	})

	t.Run("with OAuth provider", func(t *testing.T) {
		provider := "google"
		providerID := "123456"
		tokenExp := int64(1234567890)

		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID.String(),
			},
			UserID:           userID,
			Email:            email,
			Provider:         provider,
			ProviderID:       providerID,
			ProviderTokenExp: &tokenExp,
			AuthMethod:       AuthMethodOAuth,
		}

		if claims.Provider != provider {
			t.Errorf("Provider = %v, want %v", claims.Provider, provider)
		}
		if claims.ProviderID != providerID {
			t.Errorf("ProviderID = %v, want %v", claims.ProviderID, providerID)
		}
		if claims.AuthMethod != AuthMethodOAuth {
			t.Errorf("AuthMethod = %v, want %v", claims.AuthMethod, AuthMethodOAuth)
		}
		if *claims.ProviderTokenExp != tokenExp {
			t.Errorf("ProviderTokenExp = %v, want %v", *claims.ProviderTokenExp, tokenExp)
		}
	})

	t.Run("with custom claims", func(t *testing.T) {
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: userID.String(),
			},
			UserID: userID,
			Email:  email,
			CustomClaims: map[string]interface{}{
				"tier":          "premium",
				"beta_features": true,
			},
		}

		if claims.CustomClaims["tier"] != "premium" {
			t.Errorf("CustomClaims[tier] = %v, want premium", claims.CustomClaims["tier"])
		}
		if claims.CustomClaims["beta_features"] != true {
			t.Errorf("CustomClaims[beta_features] = %v, want true", claims.CustomClaims["beta_features"])
		}
	})
}

func TestEnhancedClaims_Validation(t *testing.T) {
	t.Run("HasPermission", func(t *testing.T) {
		claims := &EnhancedClaims{
			Permissions: []string{"read", "write"},
			Organizations: []OrganizationClaim{
				{
					ID:          uuid.New(),
					Permissions: []string{"delete"},
				},
			},
		}
		claims.OrganizationID = claims.Organizations[0].ID

		if !claims.HasPermission("read") {
			t.Error("Should have 'read' permission")
		}
		if !claims.HasPermission("write") {
			t.Error("Should have 'write' permission")
		}
		if !claims.HasPermission("delete") {
			t.Error("Should have 'delete' permission from organization")
		}
		if claims.HasPermission("admin") {
			t.Error("Should not have 'admin' permission")
		}
	})

	t.Run("HasRole", func(t *testing.T) {
		claims := &EnhancedClaims{
			Roles:            []string{"developer", "tester"},
			OrganizationRole: "admin",
		}

		if !claims.HasRole("developer") {
			t.Error("Should have 'developer' role")
		}
		if !claims.HasRole("admin") {
			t.Error("Should have 'admin' role from organization")
		}
		if claims.HasRole("manager") {
			t.Error("Should not have 'manager' role")
		}
	})

	t.Run("HasScope", func(t *testing.T) {
		claims := &EnhancedClaims{
			Scopes: []string{"read:profile", "write:posts"},
		}

		if !claims.HasScope("read:profile") {
			t.Error("Should have 'read:profile' scope")
		}
		if claims.HasScope("delete:all") {
			t.Error("Should not have 'delete:all' scope")
		}
	})

	t.Run("IsOrgMember", func(t *testing.T) {
		orgID1 := uuid.New()
		orgID2 := uuid.New()
		orgID3 := uuid.New()

		claims := &EnhancedClaims{
			OrganizationID: orgID1,
			Organizations: []OrganizationClaim{
				{ID: orgID2, Name: "Org 2"},
			},
		}

		if !claims.IsOrgMember(orgID1) {
			t.Error("Should be member of primary organization")
		}
		if !claims.IsOrgMember(orgID2) {
			t.Error("Should be member of secondary organization")
		}
		if claims.IsOrgMember(orgID3) {
			t.Error("Should not be member of unknown organization")
		}
	})

	t.Run("GetOrgRole", func(t *testing.T) {
		orgID1 := uuid.New()
		orgID2 := uuid.New()

		claims := &EnhancedClaims{
			OrganizationID:   orgID1,
			OrganizationRole: "owner",
			Organizations: []OrganizationClaim{
				{ID: orgID2, Role: "member"},
			},
		}

		if role := claims.GetOrgRole(orgID1); role != "owner" {
			t.Errorf("GetOrgRole(orgID1) = %v, want owner", role)
		}
		if role := claims.GetOrgRole(orgID2); role != "member" {
			t.Errorf("GetOrgRole(orgID2) = %v, want member", role)
		}
		if role := claims.GetOrgRole(uuid.New()); role != "" {
			t.Errorf("GetOrgRole(unknown) = %v, want empty", role)
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Run("valid claims", func(t *testing.T) {
			claims := &EnhancedClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
				},
				UserID: uuid.New(),
				Email:  "test@example.com",
			}

			if !claims.IsValid() {
				t.Error("Claims should be valid")
			}
		})

		t.Run("expired claims", func(t *testing.T) {
			claims := &EnhancedClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				},
				UserID: uuid.New(),
				Email:  "test@example.com",
			}

			if claims.IsValid() {
				t.Error("Expired claims should not be valid")
			}
		})

		t.Run("not yet valid", func(t *testing.T) {
			claims := &EnhancedClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				},
				UserID: uuid.New(),
				Email:  "test@example.com",
			}

			if claims.IsValid() {
				t.Error("Future claims should not be valid")
			}
		})

		t.Run("missing user ID", func(t *testing.T) {
			claims := &EnhancedClaims{
				Email: "test@example.com",
			}

			if claims.IsValid() {
				t.Error("Claims without UserID should not be valid")
			}
		})

		t.Run("missing email", func(t *testing.T) {
			claims := &EnhancedClaims{
				UserID: uuid.New(),
			}

			if claims.IsValid() {
				t.Error("Claims without Email should not be valid")
			}
		})
	})

	t.Run("ValidateForEndpoint", func(t *testing.T) {
		claims := &EnhancedClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
			UserID:      uuid.New(),
			Email:       "test@example.com",
			Scopes:      []string{"read:profile", "write:posts"},
			Permissions: []string{"post.create", "post.edit"},
		}

		t.Run("valid for endpoint", func(t *testing.T) {
			if !claims.ValidateForEndpoint(
				[]string{"read:profile"},
				[]string{"post.create"},
			) {
				t.Error("Should be valid for endpoint")
			}
		})

		t.Run("missing scope", func(t *testing.T) {
			if claims.ValidateForEndpoint(
				[]string{"admin:all"},
				[]string{"post.create"},
			) {
				t.Error("Should not be valid without required scope")
			}
		})

		t.Run("missing permission", func(t *testing.T) {
			if claims.ValidateForEndpoint(
				[]string{"read:profile"},
				[]string{"user.delete"},
			) {
				t.Error("Should not be valid without required permission")
			}
		})
	})
}

func TestRefreshClaims(t *testing.T) {
	t.Run("basic refresh claims", func(t *testing.T) {
		userID := uuid.New()
		sessionID := "session-123"

		claims := &RefreshClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
			UserID:     userID,
			TokenType:  RefreshTokenType,
			SessionID:  sessionID,
			AuthMethod: AuthMethodPassword,
		}

		if claims.UserID != userID {
			t.Errorf("UserID = %v, want %v", claims.UserID, userID)
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

func TestAPIKeyClaims(t *testing.T) {
	t.Run("API key claims", func(t *testing.T) {
		keyID := "key-123"
		userID := uuid.New()
		orgID := uuid.New()

		claims := &APIKeyClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			KeyID:          keyID,
			UserID:         userID,
			OrganizationID: orgID,
			Name:           "Production API Key",
			Scopes:         []string{"read:data", "write:data"},
			RateLimit:      1000,
		}

		if claims.KeyID != keyID {
			t.Errorf("KeyID = %v, want %v", claims.KeyID, keyID)
		}
		if claims.UserID != userID {
			t.Errorf("UserID = %v, want %v", claims.UserID, userID)
		}
		if claims.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", claims.OrganizationID, orgID)
		}
		if len(claims.Scopes) != 2 {
			t.Errorf("Scopes length = %v, want 2", len(claims.Scopes))
		}
		if claims.RateLimit != 1000 {
			t.Errorf("RateLimit = %v, want 1000", claims.RateLimit)
		}
	})
}
