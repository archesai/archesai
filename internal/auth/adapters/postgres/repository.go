// Package postgres provides PostgreSQL repository implementations
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/archesai/archesai/internal/auth/adapters"
	"github.com/archesai/archesai/internal/auth/domain"
	"github.com/archesai/archesai/internal/storage/database/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Compile-time check that Repository implements the domain.Repository interface
var _ domain.Repository = (*Repository)(nil)

// Repository handles auth data persistence using PostgreSQL
// Note: Currently uses existing generated types which may not include all auth fields
// The password_hash field needs to be added to the schema and queries
type Repository struct {
	q postgresql.Querier
}

// NewPostgresRepository creates a new auth repository
func NewPostgresRepository(q postgresql.Querier) *Repository {
	return &Repository{q: q}
}

// User operations

// GetUserByEmail retrieves a user by their email address.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToDomain(&row), nil
}

// GetUserByID retrieves a user by their unique identifier.
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	// Note: Generated queries expect string IDs, need to convert
	row, err := r.q.GetUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToDomain(&row), nil
}

// CreateUser creates a new user in the database.
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	// Note: Current schema doesn't include password_hash field
	// This needs to be added to the database schema and queries
	var imagePtr *string
	if user.Image != "" {
		imagePtr = &user.Image
	}

	params := postgresql.CreateUserParams{
		Email:         string(user.Email),
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		Image:         imagePtr,
	}

	_, err := r.q.CreateUser(ctx, params)
	// TODO: Store password hash separately or add to schema
	return err
}

// UpdateUser updates an existing user's information.
func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	email := string(user.Email)
	var imagePtr *string
	if user.Image != "" {
		imagePtr = &user.Image
	}

	params := postgresql.UpdateUserParams{
		Id:            user.Id.String(),
		Email:         &email,
		Name:          &user.Name,
		EmailVerified: pgtype.Bool{Bool: user.EmailVerified, Valid: true},
		Image:         imagePtr,
	}

	_, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
	}
	return err
}

// DeleteUser removes a user from the database.
func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
	}
	return err
}

// ListUsers retrieves a paginated list of users.
func (r *Repository) ListUsers(ctx context.Context, limit, offset int32) ([]*domain.User, error) {
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

	users := make([]*domain.User, len(rows))
	for i, row := range rows {
		users[i] = r.dbUserToDomain(&row)
	}

	return users, nil
}

// Session operations

// CreateSession creates a new user session.
func (r *Repository) CreateSession(ctx context.Context, session *domain.Session) error {
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
		ExpiresAt:            expiresAt,
		ActiveOrganizationId: activeOrgID,
		IpAddress:            ipAddress,
		UserAgent:            userAgent,
	}

	_, err := r.q.CreateSession(ctx, params)
	return err
}

// GetSessionByToken retrieves a session by its token.
func (r *Repository) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	row, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToDomain(&row), nil
}

// GetSessionByID retrieves a session by its unique identifier.
func (r *Repository) GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	row, err := r.q.GetSession(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToDomain(&row), nil
}

// UpdateSession updates an existing session.
func (r *Repository) UpdateSession(ctx context.Context, session *domain.Session) error {
	// Parse ExpiresAt string to time
	expiresAt, _ := time.Parse(time.RFC3339, session.ExpiresAt)

	var activeOrgID *string
	if session.ActiveOrganizationId != "" {
		activeOrgID = &session.ActiveOrganizationId
	}

	params := postgresql.UpdateSessionParams{
		Id:                   session.Id.String(),
		ExpiresAt:            pgtype.Timestamptz{Time: expiresAt, Valid: true},
		ActiveOrganizationId: activeOrgID,
	}

	_, err := r.q.UpdateSession(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrSessionNotFound
		}
	}
	return err
}

// DeleteSession removes a session from the database.
func (r *Repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteSession(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrSessionNotFound
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

func (r *Repository) dbUserToDomain(dbUser *postgresql.User) *domain.User {
	apiUser := adapters.AuthUserDBToAPI(dbUser)
	user := &domain.User{
		UserEntity: apiUser,
	}
	// TODO: Add password hash from account table when available
	return user
}

func (r *Repository) dbSessionToDomain(dbSession *postgresql.Session) *domain.Session {
	apiSession := adapters.AuthSessionDBToAPI(dbSession)
	session := &domain.Session{
		SessionEntity: apiSession,
	}
	return session
}
