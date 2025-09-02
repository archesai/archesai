// Package postgres provides PostgreSQL implementations for auth domain repositories.
package postgres

import (
	"context"
	"errors"

	"github.com/archesai/archesai/internal/domains/auth/entities"
	"github.com/archesai/archesai/internal/domains/auth/repositories"
	"github.com/archesai/archesai/internal/generated/database/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Compile-time check: Repository implements the port
var _ repositories.Repository = (*Repository)(nil)

// Repository handles auth data persistence using PostgreSQL
// Note: Currently uses existing generated types which may not include all auth fields
// The password_hash field needs to be added to the schema and queries
type Repository struct {
	q postgresql.Querier
}

// NewRepository creates a new auth repository
func NewRepository(q postgresql.Querier) *Repository {
	return &Repository{q: q}
}

// User operations

// GetUserByEmail retrieves a user by their email address.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToDomain(row), nil
}

// GetUserByID retrieves a user by their unique identifier.
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	// Note: Generated queries expect string IDs, need to convert
	row, err := r.q.GetUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToDomain(row), nil
}

// CreateUser creates a new user in the database.
func (r *Repository) CreateUser(ctx context.Context, user *entities.User) error {
	// Note: Current schema doesn't include password_hash field
	// This needs to be added to the database schema and queries
	params := postgresql.CreateUserParams{
		Email:         user.Email,
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		Image:         user.Image,
	}

	_, err := r.q.CreateUser(ctx, params)
	// TODO: Store password hash separately or add to schema
	return err
}

// UpdateUser updates an existing user's information.
func (r *Repository) UpdateUser(ctx context.Context, user *entities.User) error {
	params := postgresql.UpdateUserParams{
		ID:            user.ID.String(),
		Email:         &user.Email,
		Name:          &user.Name,
		EmailVerified: pgtype.Bool{Bool: user.EmailVerified, Valid: true},
		Image:         user.Image,
	}

	_, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.ErrUserNotFound
		}
	}
	return err
}

// DeleteUser removes a user from the database.
func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.ErrUserNotFound
		}
	}
	return err
}

// ListUsers retrieves a paginated list of users.
func (r *Repository) ListUsers(ctx context.Context, limit, offset int32) ([]*entities.User, error) {
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

	users := make([]*entities.User, len(rows))
	for i, row := range rows {
		users[i] = r.dbUserToDomain(row)
	}

	return users, nil
}

// Session operations

// CreateSession creates a new user session.
func (r *Repository) CreateSession(ctx context.Context, session *entities.Session) error {
	params := postgresql.CreateSessionParams{
		UserID:               session.UserID.String(),
		Token:                session.Token,
		ExpiresAt:            session.ExpiresAt,
		ActiveOrganizationID: convertUUIDToString(session.ActiveOrganizationID),
		IpAddress:            session.IPAddress,
		UserAgent:            session.UserAgent,
	}

	_, err := r.q.CreateSession(ctx, params)
	return err
}

// GetSessionByToken retrieves a session by its token.
func (r *Repository) GetSessionByToken(ctx context.Context, token string) (*entities.Session, error) {
	row, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToDomain(row), nil
}

// GetSessionByID retrieves a session by its unique identifier.
func (r *Repository) GetSessionByID(ctx context.Context, id uuid.UUID) (*entities.Session, error) {
	row, err := r.q.GetSession(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToDomain(row), nil
}

// UpdateSession updates an existing session.
func (r *Repository) UpdateSession(ctx context.Context, session *entities.Session) error {
	params := postgresql.UpdateSessionParams{
		ID:                   session.ID.String(),
		ExpiresAt:            pgtype.Timestamptz{Time: session.ExpiresAt, Valid: true},
		ActiveOrganizationID: convertUUIDToString(session.ActiveOrganizationID),
	}

	_, err := r.q.UpdateSession(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.ErrSessionNotFound
		}
	}
	return err
}

// DeleteSession removes a session from the database.
func (r *Repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteSession(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.ErrSessionNotFound
		}
	}
	return err
}

// DeleteUserSessions removes all sessions for a specific user.
func (r *Repository) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.q.DeleteSessionsByUser(ctx, userID.String())
}

// DeleteExpiredSessions removes all expired sessions from the database.
func (r *Repository) DeleteExpiredSessions(_ context.Context) error {
	// TODO: Add DeleteExpiredSessions query to auth.sql
	// For now, return nil (no-op)
	return nil
}

// Helper methods to convert between database and domain models

func (r *Repository) dbUserToDomain(dbUser postgresql.User) *entities.User {
	// Parse string ID to UUID
	id, _ := uuid.Parse(dbUser.ID)

	return &entities.User{
		ID:            id,
		Email:         dbUser.Email,
		Name:          dbUser.Name,
		PasswordHash:  "", // TODO: Add to schema
		EmailVerified: dbUser.EmailVerified,
		Image:         dbUser.Image,
		CreatedAt:     dbUser.CreatedAt,
		UpdatedAt:     dbUser.UpdatedAt,
	}
}

func (r *Repository) dbSessionToDomain(dbSession postgresql.Session) *entities.Session {
	// Parse string IDs to UUIDs
	id, _ := uuid.Parse(dbSession.ID)
	userID, _ := uuid.Parse(dbSession.UserID)

	var activeOrgID *uuid.UUID
	if dbSession.ActiveOrganizationID != nil {
		if parsed, err := uuid.Parse(*dbSession.ActiveOrganizationID); err == nil {
			activeOrgID = &parsed
		}
	}

	return &entities.Session{
		ID:                   id,
		UserID:               userID,
		Token:                dbSession.Token,
		ActiveOrganizationID: activeOrgID,
		IPAddress:            dbSession.IpAddress,
		UserAgent:            dbSession.UserAgent,
		ExpiresAt:            dbSession.ExpiresAt,
		CreatedAt:            dbSession.CreatedAt,
		UpdatedAt:            dbSession.UpdatedAt,
	}
}

// Helper to convert UUID pointer to string pointer
func convertUUIDToString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}
