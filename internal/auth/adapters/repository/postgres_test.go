package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/auth"
	"github.com/archesai/archesai/internal/testutil"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func TestPostgresRepository_SessionOperations(t *testing.T) {
	t.Skip("Skipping test - repository methods not yet implemented")
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pgContainer := testutil.StartPostgresContainer(ctx, t)

	// Run migrations
	err := pgContainer.RunMigrations("../../../migrations/postgresql")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Open sql.DB connection for repository
	db, err := pgxpool.New(ctx, pgContainer.DSN)
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Create a user for sessions directly with SQL
	userID := uuid.New()
	_, err = db.Exec(ctx, `
		INSERT INTO "user" (id, email, name, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, "session@example.com", "Session User", false, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user := &users.User{
		Id:            userID,
		Email:         openapi_types.Email("session@example.com"),
		Name:          "Session User",
		EmailVerified: false,
	}

	t.Run("CreateSession", func(t *testing.T) {
		session := &auth.Session{
			UserId:               user.Id,
			Token:                "test-token",
			ExpiresAt:            time.Now().Add(time.Hour).Format(time.RFC3339),
			ActiveOrganizationId: uuid.Nil,
			IpAddress:            "192.168.1.1",
			UserAgent:            "TestAgent/1.0",
		}

		created, err := repo.CreateSession(ctx, session)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if created.Id == uuid.Nil {
			t.Error("Expected valid session ID, got nil")
		}
		if created.Token != session.Token {
			t.Errorf("Expected token %v, got %v", session.Token, created.Token)
		}
	})

	t.Run("GetSessionByToken", func(t *testing.T) {
		token := "unique-token"
		session := &auth.Session{
			UserId:               user.Id,
			Token:                token,
			ExpiresAt:            time.Now().Add(time.Hour).Format(time.RFC3339),
			ActiveOrganizationId: uuid.Nil,
			IpAddress:            "192.168.1.1",
			UserAgent:            "TestAgent/1.0",
		}

		_, err := repo.CreateSession(ctx, session)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		retrieved, err := repo.GetSessionByToken(ctx, token)
		if err != nil {
			t.Fatalf("Failed to get session by token: %v", err)
		}

		if retrieved.Token != token {
			t.Errorf("Expected token %s, got %s", token, retrieved.Token)
		}
	})

	t.Run("DeleteExpiredSessions", func(t *testing.T) {
		// Create an expired session
		expiredSession := &auth.Session{
			UserId:               user.Id,
			Token:                "expired-token",
			ExpiresAt:            time.Now().Add(-time.Hour).Format(time.RFC3339), // Expired
			ActiveOrganizationId: uuid.Nil,
			IpAddress:            "192.168.1.1",
			UserAgent:            "TestAgent/1.0",
		}

		_, err := repo.CreateSession(ctx, expiredSession)
		if err != nil {
			t.Fatalf("Failed to create expired session: %v", err)
		}

		// Delete expired sessions
		err = repo.DeleteExpiredSessions(ctx)
		if err != nil {
			t.Fatalf("Failed to delete expired sessions: %v", err)
		}

		// Verify the expired session is deleted
		_, err = repo.GetSessionByToken(ctx, "expired-token")
		if err == nil {
			t.Error("Expected error getting expired session, got nil")
		}
	})

	t.Run("DeleteUserSessions", func(t *testing.T) {
		// Create a new user for this test directly with SQL
		testUserID := uuid.New()
		_, err := db.Exec(ctx, `
			INSERT INTO "user" (id, email, name, email_verified, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, testUserID, "sessiondelete@example.com", "Session Delete User", false, time.Now(), time.Now())
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		testUser := &users.User{
			Id:            testUserID,
			Email:         openapi_types.Email("sessiondelete@example.com"),
			Name:          "Session Delete User",
			EmailVerified: false,
		}

		// Create multiple sessions for the user
		for i := 0; i < 3; i++ {
			session := &auth.Session{
				UserId:               testUser.Id,
				Token:                fmt.Sprintf("user-token-%d", i),
				ExpiresAt:            time.Now().Add(time.Hour).Format(time.RFC3339),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
				ActiveOrganizationId: uuid.Nil,
				IpAddress:            "192.168.1.1",
				UserAgent:            "TestAgent/1.0",
			}
			_, err := repo.CreateSession(ctx, session)
			if err != nil {
				t.Fatalf("Failed to create session %d: %v", i, err)
			}
		}

		// Delete all sessions for the user
		err = repo.DeleteUserSessions(ctx, testUser.Id)
		if err != nil {
			t.Fatalf("Failed to delete user sessions: %v", err)
		}

		// Verify all sessions are deleted
		for i := 0; i < 3; i++ {
			_, err := repo.GetSessionByToken(ctx, fmt.Sprintf("user-token-%d", i))
			if err == nil {
				t.Errorf("Expected error getting deleted session %d, got nil", i)
			}
		}
	})
}

func TestPostgresRepository_AccountOperations(t *testing.T) {
	t.Skip("Skipping test - repository methods not yet implemented")
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pgContainer := testutil.StartPostgresContainer(ctx, t)

	// Run migrations
	err := pgContainer.RunMigrations("../../../migrations/postgresql")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Open sql.DB connection for repository
	db, err := pgxpool.New(ctx, pgContainer.DSN)
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Create a user for accounts directly with SQL
	userID := uuid.New()
	_, err = db.Exec(ctx, `
		INSERT INTO "user" (id, email, name, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, "account@example.com", "Account User", false, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user := &users.User{
		Id:            userID,
		Email:         openapi_types.Email("account@example.com"),
		Name:          "Account User",
		EmailVerified: false,
	}

	t.Run("CreateAccount", func(t *testing.T) {
		account := &auth.Account{
			UserId:     user.Id,
			ProviderId: auth.Local,
			AccountId:  "account@example.com",
			Password:   "hashed-password",
		}

		created, err := repo.CreateAccount(ctx, account)
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		if created.Id == uuid.Nil {
			t.Error("Expected valid account ID, got nil")
		}
		if created.AccountId != account.AccountId {
			t.Errorf("Expected account ID %v, got %v", account.AccountId, created.AccountId)
		}
	})

	t.Run("GetAccountByProviderAndProviderID", func(t *testing.T) {
		providerID := "provider@example.com"
		account := &auth.Account{
			UserId:     user.Id,
			ProviderId: auth.Local,
			AccountId:  providerID,
			Password:   "hashed-password",
		}

		_, err := repo.CreateAccount(ctx, account)
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		retrieved, err := repo.GetAccountByProviderAndProviderID(ctx, string(auth.Local), providerID)
		if err != nil {
			t.Fatalf("Failed to get account: %v", err)
		}

		if retrieved.AccountId != providerID {
			t.Errorf("Expected account ID %s, got %s", providerID, retrieved.AccountId)
		}
	})

	t.Run("ListUserAccounts", func(t *testing.T) {
		// Create a new user for this test directly with SQL
		testUserID := uuid.New()
		_, err := db.Exec(ctx, `
			INSERT INTO "user" (id, email, name, email_verified, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, testUserID, "multiaccounts@example.com", "Multi Account User", false, time.Now(), time.Now())
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		testUser := &users.User{
			Id:            testUserID,
			Email:         openapi_types.Email("multiaccounts@example.com"),
			Name:          "Multi Account User",
			EmailVerified: false,
		}

		// Create multiple accounts for the user
		providers := []auth.AccountProviderId{auth.Local, auth.Google, auth.Github}
		for i, provider := range providers {
			account := &auth.Account{
				UserId:     testUser.Id,
				ProviderId: provider,
				AccountId:  fmt.Sprintf("account-%d@example.com", i),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			_, err := repo.CreateAccount(ctx, account)
			if err != nil {
				t.Fatalf("Failed to create account %d: %v", i, err)
			}
		}

		// Get all accounts for the user
		accounts, err := repo.ListUserAccounts(ctx, testUser.Id)
		if err != nil {
			t.Fatalf("Failed to get accounts: %v", err)
		}

		if len(accounts) != 3 {
			t.Errorf("Expected 3 accounts, got %d", len(accounts))
		}
	})
}
