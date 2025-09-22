package stores

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// MagicLinkToken represents a magic link verification token in the database.
type MagicLinkToken struct {
	ID             uuid.UUID
	UserID         *uuid.UUID
	Token          string
	TokenHash      string
	Code           string
	Identifier     string
	DeliveryMethod string
	ExpiresAt      time.Time
	UsedAt         *time.Time
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// MagicLinkRepository defines the interface for magic link persistence.
type MagicLinkRepository interface {
	Create(ctx context.Context, token *MagicLinkToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*MagicLinkToken, error)
	GetByCode(ctx context.Context, identifier, code string) (*MagicLinkToken, error)
	MarkUsed(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	CountRecentByIdentifier(ctx context.Context, identifier string, since time.Time) (int, error)
}

// PostgresMagicLinkRepository implements MagicLinkRepository for PostgreSQL.
type PostgresMagicLinkRepository struct {
	db *sql.DB
}

// NewPostgresMagicLinkRepository creates a new PostgreSQL magic link repository.
func NewPostgresMagicLinkRepository(db *sql.DB) *PostgresMagicLinkRepository {
	return &PostgresMagicLinkRepository{db: db}
}

// Create stores a new magic link token.
func (r *PostgresMagicLinkRepository) Create(ctx context.Context, token *MagicLinkToken) error {
	query := `
		INSERT INTO magic_link_tokens (
			id, user_id, token_hash, code, identifier,
			delivery_method, expires_at, ip_address, user_agent,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.Code,
		token.Identifier,
		token.DeliveryMethod,
		token.ExpiresAt,
		token.IPAddress,
		token.UserAgent,
		token.CreatedAt,
		token.UpdatedAt,
	)
	return err
}

// GetByTokenHash retrieves a token by its hash.
func (r *PostgresMagicLinkRepository) GetByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*MagicLinkToken, error) {
	var token MagicLinkToken
	query := `
		SELECT id, user_id, token_hash, code, identifier,
			   delivery_method, expires_at, used_at, ip_address,
			   user_agent, created_at, updated_at
		FROM magic_link_tokens
		WHERE token_hash = $1 AND used_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.Code,
		&token.Identifier,
		&token.DeliveryMethod,
		&token.ExpiresAt,
		&token.UsedAt,
		&token.IPAddress,
		&token.UserAgent,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrTokenNotFound
	}
	return &token, err
}

// GetByCode retrieves a token by identifier and code.
func (r *PostgresMagicLinkRepository) GetByCode(
	ctx context.Context,
	identifier, code string,
) (*MagicLinkToken, error) {
	var token MagicLinkToken
	query := `
		SELECT id, user_id, token_hash, code, identifier,
			   delivery_method, expires_at, used_at, ip_address,
			   user_agent, created_at, updated_at
		FROM magic_link_tokens
		WHERE identifier = $1 AND code = $2 AND used_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, identifier, code).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.Code,
		&token.Identifier,
		&token.DeliveryMethod,
		&token.ExpiresAt,
		&token.UsedAt,
		&token.IPAddress,
		&token.UserAgent,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrTokenNotFound
	}
	return &token, err
}

// MarkUsed marks a token as used.
func (r *PostgresMagicLinkRepository) MarkUsed(ctx context.Context, tokenID uuid.UUID) error {
	query := `UPDATE magic_link_tokens SET used_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	return err
}

// DeleteExpired removes expired tokens.
func (r *PostgresMagicLinkRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM magic_link_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

// CountRecentByIdentifier counts recent tokens for rate limiting.
func (r *PostgresMagicLinkRepository) CountRecentByIdentifier(
	ctx context.Context,
	identifier string,
	since time.Time,
) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM magic_link_tokens
		WHERE identifier = $1 AND created_at >= $2`

	err := r.db.QueryRowContext(ctx, query, identifier, since).Scan(&count)
	return count, err
}

// Errors
var (
	ErrTokenNotFound = sql.ErrNoRows
)
