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
	// Create user params with required fields
	params := postgresql.CreateUserParams{
		Email: string(entity.Email),
		Name:  entity.Name,
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

	params := postgresql.UpdateUserParams{
		Id:   id,
		Name: name,
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

	// TODO: Get actual total count
	return users, int64(len(users)), nil
}

// Session operations

// CreateSession creates a new user session.
func (r *PostgresRepository) CreateSession(ctx context.Context, entity *Session) (*Session, error) {
	// Parse ExpiresAt string to time
	expiresAt, _ := time.Parse(time.RFC3339, entity.ExpiresAt)

	// TODO: Add token generation and storage
	token := uuid.New().String()

	params := postgresql.CreateSessionParams{
		UserId:               entity.UserId,
		Token:                token,
		ActiveOrganizationId: &entity.ActiveOrganizationId,
		IpAddress:            nil, // TODO: Get from context
		UserAgent:            nil, // TODO: Get from context
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
	// Parse token as UUID for now since we're using UUID as token
	tokenUUID, err := uuid.Parse(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}
	// For now, token is the session ID
	// TODO: Implement proper token mechanism
	row, err := r.q.GetSession(ctx, tokenUUID)
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
func (r *PostgresRepository) CreateAccount(_ context.Context, _ *Account) (*Account, error) {
	// TODO: Implement when account table is added to schema
	return nil, fmt.Errorf("account operations not yet implemented")
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
	return &User{
		Id:    dbUser.Id,
		Email: openapi_types.Email(dbUser.Email),
		Name:  dbUser.Name,
		// TODO: Add other fields as needed
	}
}

func (r *PostgresRepository) dbSessionToEntity(dbSession *postgresql.Session) *Session {
	session := &Session{
		Id:        dbSession.Id,
		UserId:    dbSession.UserId,
		ExpiresAt: dbSession.ExpiresAt.Format(time.RFC3339),
	}
	if dbSession.ActiveOrganizationId != nil {
		session.ActiveOrganizationId = *dbSession.ActiveOrganizationId
	}
	return session
}

// GetUserByUsername retrieves a user by username
func (r *PostgresRepository) GetUserByUsername(_ context.Context, _ string) (*User, error) {
	// TODO: Implement when username field is added to users
	return nil, fmt.Errorf("not implemented yet - username field not in schema")
}

// GetAccountByProviderAndProviderID retrieves an account by provider and provider ID
func (r *PostgresRepository) GetAccountByProviderAndProviderID(_ context.Context, _, _ string) (*Account, error) {
	// TODO: Implement when Account queries are added to postgresql
	return nil, fmt.Errorf("not implemented yet - waiting for Account SQL queries")
}

// GetAccountsByUserID retrieves accounts by user ID
func (r *PostgresRepository) GetAccountsByUserID(_ context.Context, _ uuid.UUID) ([]*Account, error) {
	// TODO: Implement when Account queries are added to postgresql
	return nil, fmt.Errorf("not implemented yet - waiting for Account SQL queries")
}
