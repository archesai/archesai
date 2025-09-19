-- +goose Up
-- +goose StatementBegin

-- Alter existing session table to add new columns for magic link support
ALTER TABLE session
    RENAME COLUMN active_organization_id TO organization_id;

ALTER TABLE session
    ADD COLUMN IF NOT EXISTS auth_method VARCHAR(50),
    ADD COLUMN IF NOT EXISTS auth_provider VARCHAR(50),
    ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Create magic link tokens table
CREATE TABLE IF NOT EXISTS magic_link_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    code VARCHAR(6), -- Optional 6-digit OTP code
    identifier VARCHAR(255) NOT NULL, -- Email or username
    delivery_method VARCHAR(50), -- email, console, webhook, otp
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add indexes for quick lookups
CREATE INDEX idx_magic_link_tokens_token_hash ON magic_link_tokens(token_hash);
CREATE INDEX idx_magic_link_tokens_code ON magic_link_tokens(code) WHERE code IS NOT NULL;
CREATE INDEX idx_magic_link_tokens_identifier ON magic_link_tokens(identifier);
CREATE INDEX idx_magic_link_tokens_expires_at ON magic_link_tokens(expires_at) WHERE used_at IS NULL;

-- Sessions table already created with the correct columns above

-- Create function to clean up expired tokens
CREATE OR REPLACE FUNCTION cleanup_expired_magic_links() RETURNS void AS $$
BEGIN
    DELETE FROM magic_link_tokens
    WHERE expires_at < NOW()
    AND used_at IS NULL;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the cleanup function
DROP FUNCTION IF EXISTS cleanup_expired_magic_links();

-- Drop indexes
DROP INDEX IF EXISTS idx_magic_link_tokens_token_hash;
DROP INDEX IF EXISTS idx_magic_link_tokens_code;
DROP INDEX IF EXISTS idx_magic_link_tokens_identifier;
DROP INDEX IF EXISTS idx_magic_link_tokens_expires_at;

-- Drop magic link tokens table
DROP TABLE IF EXISTS magic_link_tokens;

-- Revert session table changes
ALTER TABLE session
    RENAME COLUMN organization_id TO active_organization_id;

ALTER TABLE session
    DROP COLUMN IF EXISTS auth_method,
    DROP COLUMN IF EXISTS auth_provider,
    DROP COLUMN IF EXISTS metadata;

-- +goose StatementEnd
