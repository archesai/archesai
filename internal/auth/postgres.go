// Package auth provides PostgreSQL repository implementations
package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// PostgresRepository handles auth data persistence using PostgreSQL
// Note: Currently uses existing generated types which may not include all auth fields
// The password_hash field needs to be added to the schema and queries
type PostgresRepository struct {
	q postgresql.Querier
}

// Compile-time check that PostgresRepository implements the Repository interface
var _ Repository = (*PostgresRepository)(nil)

// NewPostgresRepository creates a new auth repository
func NewPostgresRepository(q postgresql.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

// User operations

// GetUserByEmail retrieves a user by their email address.
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToEntity(&row), nil
}

// GetUser retrieves a user by their unique identifier.
func (r *PostgresRepository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	// Note: Generated queries expect string IDs, need to convert
	row, err := r.q.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToEntity(&row), nil
}

// CreateUser creates a new user in the database.
func (r *PostgresRepository) CreateUser(ctx context.Context, entity *User) (*User, error) {
	// Generate ID if not provided
	if entity.Id == uuid.Nil {
		entity.Id = uuid.New()
	}

	// Create user params with required fields
	params := postgresql.CreateUserParams{
		Id:            entity.Id,
		Email:         string(entity.Email),
		Name:          entity.Name,
		EmailVerified: entity.EmailVerified,
		Image:         &entity.Image,
	}

	// Create the user
	dbUser, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	// TODO: Also create account with password_hash when auth schema is ready
	// For now, we'll need to handle authentication separately

	return r.dbUserToEntity(&dbUser), nil
}

// UpdateUser updates an existing user's information.
func (r *PostgresRepository) UpdateUser(ctx context.Context, id uuid.UUID, entity *User) (*User, error) {
	// Create update params
	var name *string
	if entity.Name != "" {
		name = &entity.Name
	}

	var email *string
	if entity.Email != "" {
		emailStr := string(entity.Email)
		email = &emailStr
	}

	params := postgresql.UpdateUserParams{
		Id:            id,
		Name:          name,
		Email:         email,
		EmailVerified: &entity.EmailVerified,
		Image:         &entity.Image,
	}

	// Update the user
	dbUser, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToEntity(&dbUser), nil
}

// DeleteUser removes a user from the database.
func (r *PostgresRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
	}
	return err
}

// ListUsers retrieves a paginated list of users.
func (r *PostgresRepository) ListUsers(ctx context.Context, params ListUsersParams) ([]*User, int64, error) {
	limit := int32(50)
	offset := int32(0)
	if params.Limit > 0 {
		limit = int32(params.Limit)
	}
	if params.Offset > 0 {
		offset = int32(params.Offset)
	}

	dbParams := postgresql.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	}

	rows, err := r.q.ListUsers(ctx, dbParams)
	if err != nil {
		return nil, 0, err
	}

	users := make([]*User, len(rows))
	for i, row := range rows {
		users[i] = r.dbUserToEntity(&row)
	}

	// Get total count
	total, err := r.q.CountUsers(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Session operations

// CreateSession creates a new user session.
func (r *PostgresRepository) CreateSession(ctx context.Context, entity *Session) (*Session, error) {
	// Generate ID if not provided
	if entity.Id == uuid.Nil {
		entity.Id = uuid.New()
	}

	// Parse ExpiresAt string to time
	expiresAt, _ := time.Parse(time.RFC3339, entity.ExpiresAt)

	// Use provided token or generate a new one
	token := entity.Token
	if token == "" {
		token = uuid.New().String()
	}

	var activeOrgID *uuid.UUID
	if entity.ActiveOrganizationId != uuid.Nil {
		activeOrgID = &entity.ActiveOrganizationId
	}

	var ipAddress, userAgent *string
	if entity.IpAddress != "" {
		ipAddress = &entity.IpAddress
	}
	if entity.UserAgent != "" {
		userAgent = &entity.UserAgent
	}

	params := postgresql.CreateSessionParams{
		Id:                   entity.Id,
		UserId:               entity.UserId,
		Token:                token,
		ActiveOrganizationId: activeOrgID,
		IpAddress:            ipAddress,
		UserAgent:            userAgent,
		ExpiresAt:            expiresAt,
	}

	dbSession, err := r.q.CreateSession(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.dbSessionToEntity(&dbSession), nil
}

// GetSessionByToken retrieves a session by its token.
func (r *PostgresRepository) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	row, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToEntity(&row), nil
}

// GetSession retrieves a session by its ID.
func (r *PostgresRepository) GetSession(ctx context.Context, id uuid.UUID) (*Session, error) {
	row, err := r.q.GetSession(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToEntity(&row), nil
}

// UpdateSession updates a session's information.
func (r *PostgresRepository) UpdateSession(ctx context.Context, id uuid.UUID, entity *Session) (*Session, error) {
	// Parse ExpiresAt string to time
	expiresAt, _ := time.Parse(time.RFC3339, entity.ExpiresAt)

	params := postgresql.UpdateSessionParams{
		Id:        id,
		ExpiresAt: &expiresAt,
	}

	dbSession, err := r.q.UpdateSession(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToEntity(&dbSession), nil
}

// DeleteSession removes a session from the database.
func (r *PostgresRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteSession(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSessionNotFound
		}
	}
	return err
}

// DeleteUserSessions removes all sessions for a specific user.
func (r *PostgresRepository) DeleteUserSessions(_ context.Context, _ uuid.UUID) error {
	// TODO: Implement bulk delete query
	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func (r *PostgresRepository) DeleteExpiredSessions(_ context.Context) error {
	// TODO: Implement cleanup query
	return nil
}

// Account operations

// CreateAccount creates a new account for a user
func (r *PostgresRepository) CreateAccount(ctx context.Context, entity *Account) (*Account, error) {
	// Generate ID if not provided
	if entity.Id == uuid.Nil {
		entity.Id = uuid.New()
	}

	var password *string
	if entity.Password != "" {
		password = &entity.Password
	}

	params := postgresql.CreateAccountParams{
		Id:         entity.Id,
		UserId:     entity.UserId,
		ProviderId: string(entity.ProviderId),
		AccountId:  entity.AccountId,
		Password:   password,
	}

	dbAccount, err := r.q.CreateAccount(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.dbAccountToEntity(&dbAccount), nil
}

// GetAccount retrieves an account by its ID
func (r *PostgresRepository) GetAccount(_ context.Context, _ uuid.UUID) (*Account, error) {
	// TODO: Implement when account table is added to schema
	return nil, fmt.Errorf("account operations not yet implemented")
}

// UpdateAccount updates an existing account
func (r *PostgresRepository) UpdateAccount(_ context.Context, _ uuid.UUID, _ *Account) (*Account, error) {
	// TODO: Implement when account table is added to schema
	return nil, fmt.Errorf("account operations not yet implemented")
}

// DeleteAccount removes an account
func (r *PostgresRepository) DeleteAccount(_ context.Context, _ uuid.UUID) error {
	// TODO: Implement when account table is added to schema
	return fmt.Errorf("account operations not yet implemented")
}

// ListAccounts lists accounts with pagination
func (r *PostgresRepository) ListAccounts(_ context.Context, _ ListAccountsParams) ([]*Account, int64, error) {
	// TODO: Implement when account table is added to schema
	return nil, 0, fmt.Errorf("account operations not yet implemented")
}

// GetAccountByProviderID retrieves an account by provider and provider ID
func (r *PostgresRepository) GetAccountByProviderID(_ context.Context, _, _ string) (*Account, error) {
	// TODO: Implement when account table is added to schema
	return nil, fmt.Errorf("account operations not yet implemented")
}

// ListUserAccounts retrieves all accounts for a specific user
func (r *PostgresRepository) ListUserAccounts(_ context.Context, _ uuid.UUID) ([]*Account, error) {
	// TODO: Implement when account table is added to schema
	return nil, fmt.Errorf("account operations not yet implemented")
}

// ListSessions retrieves a paginated list of sessions.
func (r *PostgresRepository) ListSessions(_ context.Context, _ ListSessionsParams) ([]*Session, int64, error) {
	// TODO: Implement when session list query is available
	return nil, 0, fmt.Errorf("list sessions not yet implemented")
}

// DeleteSessionByToken deletes a session by its token.
func (r *PostgresRepository) DeleteSessionByToken(ctx context.Context, token string) error {
	// For now, token is the session ID
	// TODO: Implement proper token mechanism
	tokenUUID, err := uuid.Parse(token)
	if err != nil {
		return fmt.Errorf("invalid token format: %w", err)
	}
	err = r.q.DeleteSession(ctx, tokenUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSessionNotFound
		}
	}
	return err
}

// Helper methods to convert between database and entities

func (r *PostgresRepository) dbUserToEntity(dbUser *postgresql.User) *User {
	user := &User{
		Id:            dbUser.Id,
		Email:         openapi_types.Email(dbUser.Email),
		Name:          dbUser.Name,
		EmailVerified: dbUser.EmailVerified,
		CreatedAt:     dbUser.CreatedAt,
		UpdatedAt:     dbUser.UpdatedAt,
	}
	if dbUser.Image != nil {
		user.Image = *dbUser.Image
	}
	return user
}

func (r *PostgresRepository) dbSessionToEntity(dbSession *postgresql.Session) *Session {
	session := &Session{
		Id:        dbSession.Id,
		UserId:    dbSession.UserId,
		Token:     dbSession.Token,
		ExpiresAt: dbSession.ExpiresAt.Format(time.RFC3339),
		CreatedAt: dbSession.CreatedAt,
		UpdatedAt: dbSession.UpdatedAt,
	}
	if dbSession.ActiveOrganizationId != nil {
		session.ActiveOrganizationId = *dbSession.ActiveOrganizationId
	}
	if dbSession.IpAddress != nil {
		session.IpAddress = *dbSession.IpAddress
	}
	if dbSession.UserAgent != nil {
		session.UserAgent = *dbSession.UserAgent
	}
	return session
}

func (r *PostgresRepository) dbAccountToEntity(dbAccount *postgresql.Account) *Account {
	account := &Account{
		Id:         dbAccount.Id,
		UserId:     dbAccount.UserId,
		ProviderId: AccountProviderId(dbAccount.ProviderId),
		AccountId:  dbAccount.AccountId,
		CreatedAt:  dbAccount.CreatedAt,
		UpdatedAt:  dbAccount.UpdatedAt,
	}
	if dbAccount.Password != nil {
		account.Password = *dbAccount.Password
	}
	return account
}

// GetUserByUsername retrieves a user by username
func (r *PostgresRepository) GetUserByUsername(_ context.Context, _ string) (*User, error) {
	// TODO: Implement when username field is added to users
	return nil, fmt.Errorf("not implemented yet - username field not in schema")
}

// GetAccountByProviderAndProviderID retrieves an account by provider and provider ID
func (r *PostgresRepository) GetAccountByProviderAndProviderID(ctx context.Context, provider string, providerAccountID string) (*Account, error) {
	// First need to get the user by the provider account ID
	// This is a bit complex, we need to find account by provider and account_id
	rows, err := r.q.ListAccounts(ctx, postgresql.ListAccountsParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row.ProviderId == provider && row.AccountId == providerAccountID {
			return r.dbAccountToEntity(&row), nil
		}
	}

	return nil, ErrAccountNotFound
}

// GetAccountsByUserID retrieves accounts by user ID
func (r *PostgresRepository) GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*Account, error) {
	rows, err := r.q.ListAccountsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, len(rows))
	for i, row := range rows {
		accounts[i] = r.dbAccountToEntity(&row)
	}

	return accounts, nil
}
