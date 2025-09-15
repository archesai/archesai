-- +goose Up
-- +goose StatementBegin

-- Update api_token table to match our new structure
-- First, drop the old foreign key constraint to organization
ALTER TABLE api_token DROP CONSTRAINT IF EXISTS api_token_user_id_fkey;

-- Add user_id column and organization_id column
ALTER TABLE api_token
ADD COLUMN user_id_new UUID,
ADD COLUMN organization_id UUID;

-- For existing records, if any exist, we'll need to handle them
-- Since this is a breaking change, we'll assume no critical data exists
-- In production, you would need a data migration strategy

-- Drop the old user_id column (which was referencing organization)
ALTER TABLE api_token DROP COLUMN user_id;

-- Rename user_id_new to user_id
ALTER TABLE api_token RENAME COLUMN user_id_new TO user_id;

-- Add NOT NULL constraints
ALTER TABLE api_token ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE api_token ALTER COLUMN organization_id SET NOT NULL;

-- Add foreign key constraints
ALTER TABLE api_token ADD CONSTRAINT api_token_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE;
ALTER TABLE api_token ADD CONSTRAINT api_token_organization_id_fkey
    FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE;

-- Update the key column to key_hash for better security naming
ALTER TABLE api_token RENAME COLUMN key TO key_hash;

-- Rename permissions to scopes and change to array
ALTER TABLE api_token DROP COLUMN IF EXISTS permissions;
ALTER TABLE api_token ADD COLUMN scopes TEXT[] NOT NULL DEFAULT '{}';

-- Add rate_limit column (consolidate rate limiting fields)
ALTER TABLE api_token ADD COLUMN rate_limit INTEGER NOT NULL DEFAULT 60;

-- Add last_used_at column
ALTER TABLE api_token ADD COLUMN last_used_at TIMESTAMPTZ;

-- Remove unused columns
ALTER TABLE api_token DROP COLUMN IF EXISTS rate_limit_enabled;
ALTER TABLE api_token DROP COLUMN IF EXISTS rate_limit_max;
ALTER TABLE api_token DROP COLUMN IF EXISTS rate_limit_time_window;
ALTER TABLE api_token DROP COLUMN IF EXISTS refill_amount;
ALTER TABLE api_token DROP COLUMN IF EXISTS refill_interval;
ALTER TABLE api_token DROP COLUMN IF EXISTS remaining;
ALTER TABLE api_token DROP COLUMN IF EXISTS request_count;
ALTER TABLE api_token DROP COLUMN IF EXISTS last_refill;
ALTER TABLE api_token DROP COLUMN IF EXISTS last_request;
ALTER TABLE api_token DROP COLUMN IF EXISTS enabled;
ALTER TABLE api_token DROP COLUMN IF EXISTS start;
ALTER TABLE api_token DROP COLUMN IF EXISTS metadata;

-- Update indexes
DROP INDEX IF EXISTS idx_api_token_user_id;
CREATE INDEX idx_api_token_user_id ON api_token(user_id);
CREATE INDEX idx_api_token_organization_id ON api_token(organization_id);
CREATE INDEX idx_api_token_last_used_at ON api_token(last_used_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- This is a destructive migration, down migration would lose data
-- In production, you would implement a more careful rollback strategy

-- Drop new indexes
DROP INDEX IF EXISTS idx_api_token_last_used_at;
DROP INDEX IF EXISTS idx_api_token_organization_id;
DROP INDEX IF EXISTS idx_api_token_user_id;

-- Drop foreign key constraints
ALTER TABLE api_token DROP CONSTRAINT IF EXISTS api_token_organization_id_fkey;
ALTER TABLE api_token DROP CONSTRAINT IF EXISTS api_token_user_id_fkey;

-- Remove new columns
ALTER TABLE api_token DROP COLUMN IF EXISTS last_used_at;
ALTER TABLE api_token DROP COLUMN IF EXISTS rate_limit;
ALTER TABLE api_token DROP COLUMN IF EXISTS scopes;
ALTER TABLE api_token DROP COLUMN IF EXISTS organization_id;
ALTER TABLE api_token DROP COLUMN IF EXISTS user_id;

-- Rename key_hash back to key
ALTER TABLE api_token RENAME COLUMN key_hash TO key;

-- Restore old structure (basic version)
ALTER TABLE api_token ADD COLUMN user_id UUID;
ALTER TABLE api_token ADD COLUMN enabled BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE api_token ADD COLUMN permissions TEXT;
ALTER TABLE api_token ADD COLUMN rate_limit_enabled BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE api_token ADD COLUMN rate_limit_max INTEGER;
ALTER TABLE api_token ADD COLUMN rate_limit_time_window INTEGER;
ALTER TABLE api_token ADD COLUMN refill_amount INTEGER;
ALTER TABLE api_token ADD COLUMN refill_interval INTEGER;
ALTER TABLE api_token ADD COLUMN remaining INTEGER;
ALTER TABLE api_token ADD COLUMN request_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE api_token ADD COLUMN last_refill TIMESTAMPTZ;
ALTER TABLE api_token ADD COLUMN last_request TIMESTAMPTZ;
ALTER TABLE api_token ADD COLUMN start TEXT;
ALTER TABLE api_token ADD COLUMN metadata TEXT;

-- Add foreign key back to organization (original incorrect structure)
ALTER TABLE api_token ADD CONSTRAINT api_token_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES organization(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Restore old index
CREATE INDEX idx_api_token_user_id ON api_token(user_id);

-- +goose StatementEnd