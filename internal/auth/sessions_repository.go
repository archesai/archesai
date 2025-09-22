package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresSessionsRepository implements SessionsRepository for PostgreSQL
type PostgresSessionsRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresSessionsRepository creates a new PostgreSQL sessions repository
func NewPostgresSessionsRepository(pool *pgxpool.Pool) *PostgresSessionsRepository {
	return &PostgresSessionsRepository{pool: pool}
}

// Create creates a new session
func (r *PostgresSessionsRepository) Create(
	ctx context.Context,
	entity *SessionEntity,
) (*SessionEntity, error) {
	query := `
		INSERT INTO sessions (id, user_id, token, organization_id, auth_method, auth_provider,
		                      ip_address, user_agent, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING *`

	entity.ID = uuid.New()
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()

	row := r.pool.QueryRow(ctx, query,
		entity.ID, entity.UserID, entity.Token, entity.OrganizationID,
		entity.AuthMethod, entity.AuthProvider, entity.IPAddress, entity.UserAgent,
		entity.ExpiresAt, entity.CreatedAt, entity.UpdatedAt,
	)

	var created SessionEntity
	err := row.Scan(
		&created.ID, &created.UserID, &created.Token, &created.OrganizationID,
		&created.AuthMethod, &created.AuthProvider, &created.IPAddress, &created.UserAgent,
		&created.ExpiresAt, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// Get retrieves a session by ID
func (r *PostgresSessionsRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*SessionEntity, error) {
	query := `SELECT * FROM sessions WHERE id = $1`

	var session SessionEntity
	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&session.ID, &session.UserID, &session.Token, &session.OrganizationID,
		&session.AuthMethod, &session.AuthProvider, &session.IPAddress, &session.UserAgent,
		&session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// Update updates a session
func (r *PostgresSessionsRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	entity *SessionEntity,
) (*SessionEntity, error) {
	query := `
		UPDATE sessions
		SET updated_at = $2, expires_at = $3
		WHERE id = $1
		RETURNING *`

	entity.UpdatedAt = time.Now()

	var updated SessionEntity
	row := r.pool.QueryRow(ctx, query, id, entity.UpdatedAt, entity.ExpiresAt)
	err := row.Scan(
		&updated.ID, &updated.UserID, &updated.Token, &updated.OrganizationID,
		&updated.AuthMethod, &updated.AuthProvider, &updated.IPAddress, &updated.UserAgent,
		&updated.ExpiresAt, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

// Delete deletes a session
func (r *PostgresSessionsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

// List lists sessions
func (r *PostgresSessionsRepository) List(
	ctx context.Context,
	params ListSessionsParams,
) ([]*SessionEntity, int64, error) {
	// Simplified implementation
	query := `SELECT * FROM sessions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	countQuery := `SELECT COUNT(*) FROM sessions WHERE user_id = $1`

	limit := 10
	offset := 0
	if params.Page.Limit != nil {
		limit = *params.Page.Limit
	}
	if params.Page.Offset != nil {
		offset = *params.Page.Offset
	}

	// Get count
	var count int64
	err := r.pool.QueryRow(ctx, countQuery, params.UserID).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Get sessions
	rows, err := r.pool.Query(ctx, query, params.UserID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []*SessionEntity
	for rows.Next() {
		var session SessionEntity
		err := rows.Scan(
			&session.ID, &session.UserID, &session.Token, &session.OrganizationID,
			&session.AuthMethod, &session.AuthProvider, &session.IPAddress, &session.UserAgent,
			&session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, count, nil
}

// GetByToken retrieves a session by token
func (r *PostgresSessionsRepository) GetByToken(
	ctx context.Context,
	token string,
) (*SessionEntity, error) {
	query := `SELECT * FROM sessions WHERE token = $1`

	var session SessionEntity
	row := r.pool.QueryRow(ctx, query, token)
	err := row.Scan(
		&session.ID, &session.UserID, &session.Token, &session.OrganizationID,
		&session.AuthMethod, &session.AuthProvider, &session.IPAddress, &session.UserAgent,
		&session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteByToken deletes a session by token
func (r *PostgresSessionsRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}

// DeleteByUser deletes all sessions for a user
func (r *PostgresSessionsRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

// DeleteExpired deletes expired sessions
func (r *PostgresSessionsRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
	return err
}

// SQLiteSessionsRepository implements SessionsRepository for SQLite
type SQLiteSessionsRepository struct {
	db *sql.DB
}

// NewSQLiteSessionsRepository creates a new SQLite sessions repository
func NewSQLiteSessionsRepository(db *sql.DB) *SQLiteSessionsRepository {
	return &SQLiteSessionsRepository{db: db}
}

// Create creates a new session (SQLite implementation)
func (r *SQLiteSessionsRepository) Create(
	_ context.Context,
	entity *SessionEntity,
) (*SessionEntity, error) {
	// Similar implementation but for SQLite
	// This is a stub - implement based on your SQLite schema
	return entity, nil
}

// Get retrieves a session by ID (SQLite implementation)
func (r *SQLiteSessionsRepository) Get(_ context.Context, _ uuid.UUID) (*SessionEntity, error) {
	// Stub implementation
	return &SessionEntity{}, nil
}

// Update updates a session (SQLite implementation)
func (r *SQLiteSessionsRepository) Update(
	_ context.Context,
	_ uuid.UUID,
	entity *SessionEntity,
) (*SessionEntity, error) {
	// Stub implementation
	return entity, nil
}

// Delete deletes a session (SQLite implementation)
func (r *SQLiteSessionsRepository) Delete(_ context.Context, _ uuid.UUID) error {
	// Stub implementation
	return nil
}

// List lists sessions (SQLite implementation)
func (r *SQLiteSessionsRepository) List(
	_ context.Context,
	_ ListSessionsParams,
) ([]*SessionEntity, int64, error) {
	// Stub implementation
	return []*SessionEntity{}, 0, nil
}

// GetByToken retrieves a session by token (SQLite implementation)
func (r *SQLiteSessionsRepository) GetByToken(
	_ context.Context,
	_ string,
) (*SessionEntity, error) {
	// Stub implementation
	return &SessionEntity{}, nil
}

// DeleteByToken deletes a session by token (SQLite implementation)
func (r *SQLiteSessionsRepository) DeleteByToken(_ context.Context, _ string) error {
	// Stub implementation
	return nil
}

// DeleteByUser deletes all sessions for a user (SQLite implementation)
func (r *SQLiteSessionsRepository) DeleteByUser(_ context.Context, _ uuid.UUID) error {
	// Stub implementation
	return nil
}

// DeleteExpired deletes expired sessions (SQLite implementation)
func (r *SQLiteSessionsRepository) DeleteExpired(_ context.Context) error {
	// Stub implementation
	return nil
}
