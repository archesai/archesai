package magiclink

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostgresRepository implements Repository for PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create stores a new magic link token
func (r *PostgresRepository) Create(ctx context.Context, token *Token) error {
	query := `
		INSERT INTO magic_link_tokens (
			id, user_id, token_hash, code, identifier,
			delivery_method, expires_at, ip_address, user_agent, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.Code,
		token.Identifier,
		string(token.DeliveryMethod),
		token.ExpiresAt,
		token.IPAddress,
		token.UserAgent,
		token.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("inserting token: %w", err)
	}

	return nil
}

// GetByToken retrieves a token by its hash
func (r *PostgresRepository) GetByToken(ctx context.Context, tokenHash string) (*Token, error) {
	query := `
		SELECT
			id, user_id, token_hash, code, identifier,
			delivery_method, expires_at, used_at, ip_address,
			user_agent, created_at
		FROM magic_link_tokens
		WHERE token_hash = $1
	`

	var token Token
	var deliveryMethod string

	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.Code,
		&token.Identifier,
		&deliveryMethod,
		&token.ExpiresAt,
		&token.UsedAt,
		&token.IPAddress,
		&token.UserAgent,
		&token.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("querying token: %w", err)
	}

	token.DeliveryMethod = DeliveryMethod(deliveryMethod)
	return &token, nil
}

// GetByCode retrieves a token by identifier and OTP code
func (r *PostgresRepository) GetByCode(
	ctx context.Context,
	identifier string,
	code string,
) (*Token, error) {
	query := `
		SELECT
			id, user_id, token_hash, code, identifier,
			delivery_method, expires_at, used_at, ip_address,
			user_agent, created_at
		FROM magic_link_tokens
		WHERE identifier = $1 AND code = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var token Token
	var deliveryMethod string

	err := r.db.QueryRowContext(ctx, query, identifier, code).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.Code,
		&token.Identifier,
		&deliveryMethod,
		&token.ExpiresAt,
		&token.UsedAt,
		&token.IPAddress,
		&token.UserAgent,
		&token.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("querying token by code: %w", err)
	}

	token.DeliveryMethod = DeliveryMethod(deliveryMethod)
	return &token, nil
}

// MarkUsed marks a token as used
func (r *PostgresRepository) MarkUsed(ctx context.Context, tokenID uuid.UUID) error {
	query := `
		UPDATE magic_link_tokens
		SET used_at = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), tokenID)
	if err != nil {
		return fmt.Errorf("updating token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTokenNotFound
	}

	return nil
}

// DeleteExpired removes expired tokens
func (r *PostgresRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM magic_link_tokens
		WHERE expires_at < $1
		OR (used_at IS NOT NULL AND used_at < $2)
	`

	// Delete tokens that are expired or were used more than 24 hours ago
	now := time.Now()
	dayAgo := now.Add(-24 * time.Hour)

	_, err := r.db.ExecContext(ctx, query, now, dayAgo)
	if err != nil {
		return fmt.Errorf("deleting expired tokens: %w", err)
	}

	return nil
}

// CountRecentByIdentifier counts recent requests from an identifier
func (r *PostgresRepository) CountRecentByIdentifier(
	ctx context.Context,
	identifier string,
	since time.Time,
) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM magic_link_tokens
		WHERE identifier = $1 AND created_at > $2
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, identifier, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting recent tokens: %w", err)
	}

	return count, nil
}
