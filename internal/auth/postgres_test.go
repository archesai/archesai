package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/testutil"
	"github.com/google/uuid"
)

func TestPostgresRepository_UserOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pgContainer := testutil.StartPostgresContainer(ctx, t)

	// Run migrations
	err := pgContainer.RunMigrations("../migrations/postgresql")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	queries := postgresql.New(pgContainer.Pool)
	repo := NewPostgresRepository(queries)

	t.Run("CreateUser", func(t *testing.T) {
		user := &User{
			Email:         "test@example.com",
			Name:          "Test User",
			EmailVerified: false,
		}

		created, err := repo.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		if created.Id == uuid.Nil {
			t.Errorf("Expected valid user ID, got %v", created.Id)
		}
		if created.Email != user.Email {
			t.Errorf("Expected email %v, got %v", user.Email, created.Email)
		}
		if created.Name != user.Name {
			t.Errorf("Expected name %v, got %v", user.Name, created.Name)
		}
	})

	t.Run("GetUserByID", func(t *testing.T) {
		// Create a user first
		user := &User{
			Email:         "getbyid@example.com",
			Name:          "GetByID User",
			EmailVerified: false,
		}

		created, err := repo.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Now get the user
		retrieved, err := repo.GetUser(ctx, created.Id)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		if retrieved.Id != created.Id {
			t.Errorf("Expected user ID %v, got %v", created.Id, retrieved.Id)
		}
		if retrieved.Email != created.Email {
			t.Errorf("Expected email %s, got %s", created.Email, retrieved.Email)
		}
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		// Create a user first
		email := "getbyemail@example.com"
		user := &User{
			Email:         Email(email),
			Name:          "GetByEmail User",
			EmailVerified: false,
		}

		created, err := repo.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Now get the user by email
		retrieved, err := repo.GetUserByEmail(ctx, email)
		if err != nil {
			t.Fatalf("Failed to get user by email: %v", err)
		}

		if retrieved.Id != created.Id {
			t.Errorf("Expected user ID %v, got %v", created.Id, retrieved.Id)
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Create a user first
		user := &User{
			Email:         "update@example.com",
			Name:          "Update User",
			EmailVerified: false,
		}

		created, err := repo.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Update the user
		created.Name = "Updated User"
		created.EmailVerified = true
		updated, err := repo.UpdateUser(ctx, created.Id, created)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		if updated.Name != "Updated User" {
			t.Errorf("Expected name 'Updated User', got %s", updated.Name)
		}
		if !updated.EmailVerified {
			t.Error("Expected email to be verified")
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		// Create a user first
		user := &User{
			Email:         "delete@example.com",
			Name:          "Delete User",
			EmailVerified: false,
		}

		created, err := repo.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Delete the user
		err = repo.DeleteUser(ctx, created.Id)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Verify the user is deleted
		_, err = repo.GetUser(ctx, created.Id)
		if err == nil {
			t.Error("Expected error getting deleted user, got nil")
		}
	})

	t.Run("ListUsers", func(t *testing.T) {
		// Create multiple users
		for i := 0; i < 5; i++ {
			user := &User{
				Email:         Email(fmt.Sprintf("list%d@example.com", i)),
				Name:          fmt.Sprintf("List User %d", i),
				EmailVerified: i%2 == 0,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			_, err := repo.CreateUser(ctx, user)
			if err != nil {
				t.Fatalf("Failed to create user %d: %v", i, err)
			}
		}

		// List users with pagination
		params := ListUsersParams{
			Limit:  3,
			Offset: 0,
		}
		users, total, err := repo.ListUsers(ctx, params)
		if err != nil {
			t.Fatalf("Failed to list users: %v", err)
		}

		if len(users) > 3 {
			t.Errorf("Expected at most 3 users, got %d", len(users))
		}
		if total < 5 {
			t.Errorf("Expected at least 5 total users, got %d", total)
		}
	})
}

func TestPostgresRepository_SessionOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pgContainer := testutil.StartPostgresContainer(ctx, t)

	// Run migrations
	err := pgContainer.RunMigrations("../migrations/postgresql")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	queries := postgresql.New(pgContainer.Pool)
	repo := NewPostgresRepository(queries)

	// Create a user for sessions
	user := &User{
		Email:         "session@example.com",
		Name:          "Session User",
		EmailVerified: false,
	}
	createdUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	user = createdUser

	t.Run("CreateSession", func(t *testing.T) {
		session := &Session{
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
		session := &Session{
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
		expiredSession := &Session{
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
		// Create a new user for this test
		testUser := &User{
			Email:         "sessiondelete@example.com",
			Name:          "Session Delete User",
			EmailVerified: false,
		}
		createdTestUser, err := repo.CreateUser(ctx, testUser)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		testUser = createdTestUser

		// Create multiple sessions for the user
		for i := 0; i < 3; i++ {
			session := &Session{
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
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pgContainer := testutil.StartPostgresContainer(ctx, t)

	// Run migrations
	err := pgContainer.RunMigrations("../migrations/postgresql")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	queries := postgresql.New(pgContainer.Pool)
	repo := NewPostgresRepository(queries)

	// Create a user for accounts
	user := &User{
		Email:         "account@example.com",
		Name:          "Account User",
		EmailVerified: false,
	}
	createdUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	user = createdUser

	t.Run("CreateAccount", func(t *testing.T) {
		account := &Account{
			UserId:     user.Id,
			ProviderId: Local,
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
		account := &Account{
			UserId:     user.Id,
			ProviderId: Local,
			AccountId:  providerID,
			Password:   "hashed-password",
		}

		_, err := repo.CreateAccount(ctx, account)
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		retrieved, err := repo.GetAccountByProviderAndProviderID(ctx, string(Local), providerID)
		if err != nil {
			t.Fatalf("Failed to get account: %v", err)
		}

		if retrieved.AccountId != providerID {
			t.Errorf("Expected account ID %s, got %s", providerID, retrieved.AccountId)
		}
	})

	t.Run("GetAccountsByUserID", func(t *testing.T) {
		// Create a new user for this test
		testUser := &User{
			Email:         "multiaccounts@example.com",
			Name:          "Multi Account User",
			EmailVerified: false,
		}
		createdTestUser, err := repo.CreateUser(ctx, testUser)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		testUser = createdTestUser

		// Create multiple accounts for the user
		providers := []AccountProviderId{Local, Google, Github}
		for i, provider := range providers {
			account := &Account{
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
		accounts, err := repo.GetAccountsByUserID(ctx, testUser.Id)
		if err != nil {
			t.Fatalf("Failed to get accounts: %v", err)
		}

		if len(accounts) != 3 {
			t.Errorf("Expected 3 accounts, got %d", len(accounts))
		}
	})
}
