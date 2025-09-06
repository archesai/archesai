// Package postgres provides PostgreSQL repository implementations
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/archesai/archesai/internal/auth"
	"github.com/archesai/archesai/internal/auth/adapters"
	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// AuthPostgresRepository handles auth data persistence using PostgreSQL
// Note: Currently uses existing generated types which may not include all auth fields
// The password_hash field needs to be added to the schema and queries
type AuthPostgresRepository struct {
	q postgresql.Querier
}

// Compile-time check that AuthPostgresRepository implements the auth.Repository interface
var _ auth.Repository = (*AuthPostgresRepository)(nil)

// NewAuthPostgresRepository creates a new auth repository
func NewAuthPostgresRepository(q postgresql.Querier) *AuthPostgresRepository {
	return &AuthPostgresRepository{q: q}
}

// User operations

// GetUserByEmail retrieves a user by their email address.
func (r *AuthPostgresRepository) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToDomain(&row), nil
}

// GetUserByID retrieves a user by their unique identifier.
func (r *AuthPostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	// Note: Generated queries expect string IDs, need to convert
	row, err := r.q.GetUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToDomain(&row), nil
}

// CreateUser creates a new user in the database.
func (r *AuthPostgresRepository) CreateUser(ctx context.Context, user *auth.User) error {
	// Create user params with required fields
	params := postgresql.CreateUserParams{
		Email: string(user.Email),
		Name:  user.Name,
	}

	// Create the user
	_, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	// TODO: Also create account with password_hash when auth schema is ready
	// For now, we'll need to handle authentication separately

	return nil
}

// UpdateUser updates an existing user's information.
func (r *AuthPostgresRepository) UpdateUser(ctx context.Context, user *auth.User) error {
	// Create update params
	var name *string
	if user.Name != "" {
		name = &user.Name
	}

	params := postgresql.UpdateUserParams{
		Id:   user.Id.String(),
		Name: name,
	}

	// Update the user
	_, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.ErrUserNotFound
		}
	}
	return err
}

// DeleteUser removes a user from the database.
func (r *AuthPostgresRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.ErrUserNotFound
		}
	}
	return err
}

// ListUsers retrieves a paginated list of users.
func (r *AuthPostgresRepository) ListUsers(ctx context.Context, limit, offset int32) ([]*auth.User, error) {
	if limit == 0 {
		limit = 50
	}

	params := postgresql.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	}

	rows, err := r.q.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]*auth.User, len(rows))
	for i, row := range rows {
		users[i] = r.dbUserToDomain(&row)
	}

	return users, nil
}

// Session operations

// CreateSession creates a new user session.
func (r *AuthPostgresRepository) CreateSession(ctx context.Context, session *auth.Session) error {
	// Parse ExpiresAt string to time
	expiresAt, _ := time.Parse(time.RFC3339, session.ExpiresAt)

	var activeOrgID, ipAddress, userAgent *string
	if session.ActiveOrganizationId != "" {
		activeOrgID = &session.ActiveOrganizationId
	}
	if session.IpAddress != "" {
		ipAddress = &session.IpAddress
	}
	if session.UserAgent != "" {
		userAgent = &session.UserAgent
	}

	params := postgresql.CreateSessionParams{
		UserId:               session.UserId,
		Token:                session.Token,
		ActiveOrganizationId: activeOrgID,
		IpAddress:            ipAddress,
		UserAgent:            userAgent,
		ExpiresAt:            expiresAt,
	}

	_, err := r.q.CreateSession(ctx, params)
	return err
}

// GetSessionByToken retrieves a session by its token.
func (r *AuthPostgresRepository) GetSessionByToken(ctx context.Context, token string) (*auth.Session, error) {
	// For now, token is the session ID
	// TODO: Implement proper token mechanism
	row, err := r.q.GetSession(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToDomain(&row), nil
}

// GetSessionByID retrieves a session by its ID.
func (r *AuthPostgresRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*auth.Session, error) {
	row, err := r.q.GetSession(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToDomain(&row), nil
}

// UpdateSession updates a session's information.
func (r *AuthPostgresRepository) UpdateSession(ctx context.Context, session *auth.Session) error {
	// Parse ExpiresAt string to time
	expiresAt, _ := time.Parse(time.RFC3339, session.ExpiresAt)

	params := postgresql.UpdateSessionParams{
		Id:        session.Id.String(),
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	}

	_, err := r.q.UpdateSession(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.ErrSessionNotFound
		}
	}
	return err
}

// DeleteSession removes a session from the database.
func (r *AuthPostgresRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteSession(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.ErrSessionNotFound
		}
	}
	return err
}

// DeleteUserSessions removes all sessions for a specific user.
func (r *AuthPostgresRepository) DeleteUserSessions(_ context.Context, _ uuid.UUID) error {
	// TODO: Implement bulk delete query
	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func (r *AuthPostgresRepository) DeleteExpiredSessions(_ context.Context) error {
	// TODO: Implement cleanup query
	return nil
}

// Helper methods to convert between database and auth models

func (r *AuthPostgresRepository) dbUserToDomain(dbUser *postgresql.User) *auth.User {
	apiUser := adapters.AuthUserDBToAPI(dbUser)
	user := &auth.User{
		UserEntity: apiUser,
	}
	// TODO: Add password hash from account table when available
	return user
}

func (r *AuthPostgresRepository) dbSessionToDomain(dbSession *postgresql.Session) *auth.Session {
	apiSession := adapters.AuthSessionDBToAPI(dbSession)
	session := &auth.Session{
		SessionEntity: apiSession,
	}
	return session
}
