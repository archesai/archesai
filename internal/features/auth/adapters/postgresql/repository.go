package postgresql

import (
	"context"
	"errors"

	"github.com/archesai/archesai/gen/db/postgresql"
	"github.com/archesai/archesai/internal/features/auth/domain"
	"github.com/archesai/archesai/internal/features/auth/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Compile-time check: Repository implements the port
var _ ports.Repository = (*Repository)(nil)

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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToDomain(row), nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	// Note: Generated queries expect string IDs, need to convert
	row, err := r.q.GetUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.dbUserToDomain(row), nil
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
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

func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
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
			return domain.ErrUserNotFound
		}
	}
	return err
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
	}
	return err
}

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
		users[i] = r.dbUserToDomain(row)
	}

	return users, nil
}

// Session operations

func (r *Repository) CreateSession(ctx context.Context, session *domain.Session) error {
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

func (r *Repository) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	row, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToDomain(row), nil
}

func (r *Repository) GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	row, err := r.q.GetSession(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}
	return r.dbSessionToDomain(row), nil
}

func (r *Repository) UpdateSession(ctx context.Context, session *domain.Session) error {
	params := postgresql.UpdateSessionParams{
		ID:                   session.ID.String(),
		ExpiresAt:            pgtype.Timestamptz{Time: session.ExpiresAt, Valid: true},
		ActiveOrganizationID: convertUUIDToString(session.ActiveOrganizationID),
	}

	_, err := r.q.UpdateSession(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrSessionNotFound
		}
	}
	return err
}

func (r *Repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteSession(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrSessionNotFound
		}
	}
	return err
}

func (r *Repository) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.q.DeleteSessionsByUser(ctx, userID.String())
}

func (r *Repository) DeleteExpiredSessions(ctx context.Context) error {
	// TODO: Add DeleteExpiredSessions query to auth.sql
	// For now, return nil (no-op)
	return nil
}

// Helper methods to convert between database and domain models

func (r *Repository) dbUserToDomain(dbUser postgresql.User) *domain.User {
	// Parse string ID to UUID
	id, _ := uuid.Parse(dbUser.ID)

	return &domain.User{
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

func (r *Repository) dbSessionToDomain(dbSession postgresql.Session) *domain.Session {
	// Parse string IDs to UUIDs
	id, _ := uuid.Parse(dbSession.ID)
	userID, _ := uuid.Parse(dbSession.UserID)

	var activeOrgID *uuid.UUID
	if dbSession.ActiveOrganizationID != nil {
		if parsed, err := uuid.Parse(*dbSession.ActiveOrganizationID); err == nil {
			activeOrgID = &parsed
		}
	}

	return &domain.Session{
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
