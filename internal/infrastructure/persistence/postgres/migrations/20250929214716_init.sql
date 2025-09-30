-- +goose Up
-- +goose StatementBegin
-- PostgreSQL schema for sqlc
-- This is the source of truth for the database schema

-- Enable pgvector extension for AI embeddings
CREATE EXTENSION IF NOT EXISTS vector;

-- Create user table
CREATE TABLE "user" (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email TEXT NOT NULL UNIQUE,
  email_verified BOOLEAN NOT NULL DEFAULT false,
  image TEXT,
  name TEXT NOT NULL
);

-- Create organization table
CREATE TABLE organization (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  billing_email TEXT,
  credits INTEGER NOT NULL DEFAULT 0,
  logo TEXT,
  name TEXT NOT NULL,
  plan TEXT NOT NULL DEFAULT 'FREE' CHECK (
    plan IN (
      'BASIC',
      'FREE',
      'PREMIUM',
      'STANDARD',
      'UNLIMITED'
    )
  ),
  slug VARCHAR(50) NOT NULL,
  stripe_customer_identifier TEXT UNIQUE NOT NULL,
  CONSTRAINT organization_slug_format CHECK (slug ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$')
);

-- Create account table
CREATE TABLE account (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  access_token TEXT,
  access_token_expires_at TIMESTAMPTZ,
  account_identifier TEXT NOT NULL,
  id_token TEXT,
  password TEXT,
  provider TEXT NOT NULL,
  refresh_token TEXT,
  refresh_token_expires_at TIMESTAMPTZ,
  scope TEXT,
  user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE
);

-- Create session table
CREATE TABLE session (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  auth_method VARCHAR(50),
  auth_provider VARCHAR(50),
  expires_at TIMESTAMPTZ NOT NULL,
  ip_address TEXT,
  metadata JSONB,
  organization_id UUID,
  token TEXT NOT NULL UNIQUE,
  user_agent TEXT,
  user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE
);

-- Create member table
CREATE TABLE member (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  organization_id UUID NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
  role TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'owner', 'member')),
  user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE
);

-- Create invitation table
CREATE TABLE invitation (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  inviter_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
  organization_id UUID NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
  role TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'owner', 'member')),
  status TEXT NOT NULL DEFAULT 'pending'
);

-- Create verification_token table
CREATE TABLE verification_token (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMPTZ NOT NULL,
  identifier TEXT NOT NULL,
  value TEXT NOT NULL
);

-- Create magic link tokens table
CREATE TABLE magic_link_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(6),
  created_at TIMESTAMP DEFAULT NOW(),
  delivery_method VARCHAR(50),
  expires_at TIMESTAMP NOT NULL,
  identifier VARCHAR(255) NOT NULL,
  ip_address VARCHAR(45),
  token_hash VARCHAR(255) UNIQUE NOT NULL,
  used_at TIMESTAMP,
  user_agent TEXT,
  user_id UUID REFERENCES "user"(id) ON DELETE CASCADE
);

-- Create api_keys table
CREATE TABLE api_keys (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMPTZ,
  key_hash TEXT NOT NULL,
  last_used_at TIMESTAMPTZ,
  name TEXT,
  organization_id UUID NOT NULL,
  prefix TEXT,
  rate_limit INTEGER NOT NULL DEFAULT 60,
  scopes TEXT[] NOT NULL DEFAULT '{}',
  user_id UUID NOT NULL
);

-- Create pipeline table
CREATE TABLE pipeline (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  description TEXT,
  name TEXT,
  organization_id UUID NOT NULL REFERENCES organization (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create tool table
CREATE TABLE tool (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  description TEXT NOT NULL,
  input_mime_type TEXT NOT NULL DEFAULT 'application/octet-stream',
  name TEXT NOT NULL,
  organization_id UUID NOT NULL REFERENCES organization (id) ON DELETE CASCADE ON UPDATE CASCADE,
  output_mime_type TEXT NOT NULL DEFAULT 'application/octet-stream'
);

-- Create pipeline_step table
CREATE TABLE pipeline_step (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  pipeline_id UUID NOT NULL REFERENCES pipeline (id) ON DELETE CASCADE ON UPDATE CASCADE,
  tool_id UUID NOT NULL REFERENCES tool (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create pipeline_step_to_dependency table (junction table)
CREATE TABLE pipeline_step_to_dependency (
  pipeline_step_id UUID NOT NULL REFERENCES pipeline_step (id) ON DELETE CASCADE,
  prerequisite_id UUID NOT NULL REFERENCES pipeline_step (id) ON DELETE CASCADE,
  PRIMARY KEY (pipeline_step_id, prerequisite_id)
);

-- Create run table
CREATE TABLE run (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMPTZ,
  error TEXT,
  organization_id UUID NOT NULL REFERENCES organization (id) ON DELETE CASCADE ON UPDATE CASCADE,
  pipeline_id UUID NOT NULL REFERENCES pipeline (id) ON DELETE SET NULL ON UPDATE CASCADE,
  progress DOUBLE PRECISION NOT NULL DEFAULT 0,
  started_at TIMESTAMPTZ,
  status TEXT NOT NULL DEFAULT 'QUEUED' CHECK (
    status IN ('COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED')
  ),
  tool_id UUID NOT NULL REFERENCES tool (id) ON DELETE SET NULL ON UPDATE CASCADE
);

-- Create artifact table
CREATE TABLE artifact (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  credits INTEGER NOT NULL DEFAULT 0,
  description TEXT,
  embedding vector(1536),
  mime_type TEXT NOT NULL DEFAULT 'application/octet-stream',
  name TEXT,
  organization_id UUID NOT NULL REFERENCES organization (id) ON DELETE CASCADE ON UPDATE CASCADE,
  preview_image TEXT,
  producer_id UUID REFERENCES run (id) ON DELETE SET NULL ON UPDATE CASCADE,
  text TEXT,
  url TEXT
);

-- Create run_to_artifact table (junction table)
CREATE TABLE run_to_artifact (
  run_id UUID NOT NULL REFERENCES run (id) ON DELETE CASCADE,
  artifact_id UUID NOT NULL REFERENCES artifact (id) ON DELETE CASCADE,
  PRIMARY KEY (run_id, artifact_id)
);

-- Create label table
CREATE TABLE label (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  organization_id UUID NOT NULL REFERENCES organization (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create label_to_artifact table (junction table)
CREATE TABLE label_to_artifact (
  label_id UUID NOT NULL REFERENCES label (id) ON DELETE CASCADE,
  artifact_id UUID NOT NULL REFERENCES artifact (id) ON DELETE CASCADE,
  PRIMARY KEY (label_id, artifact_id)
);

-- Add foreign key constraints for api_keys
ALTER TABLE api_keys
ADD CONSTRAINT api_keys_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE;

ALTER TABLE api_keys
ADD CONSTRAINT api_keys_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE;

-- Create indexes
CREATE INDEX idx_account_user_id ON account (user_id);
CREATE INDEX idx_session_user_id ON session (user_id);
CREATE INDEX idx_member_organization_id ON member (organization_id);
CREATE INDEX idx_member_user_id ON member (user_id);
CREATE INDEX idx_invitation_organization_id ON invitation (organization_id);
CREATE INDEX idx_invitation_inviter_id ON invitation (inviter_id);
CREATE INDEX idx_api_keys_user_id ON api_keys (user_id);
CREATE INDEX idx_api_keys_organization_id ON api_keys (organization_id);
CREATE INDEX idx_api_keys_last_used_at ON api_keys (last_used_at);
CREATE INDEX idx_api_keys_key_hash ON api_keys (key_hash);
CREATE INDEX idx_pipeline_organization_id ON pipeline (organization_id);
CREATE INDEX idx_tool_organization_id ON tool (organization_id);
CREATE INDEX idx_pipeline_step_pipeline_id ON pipeline_step (pipeline_id);
CREATE INDEX idx_pipeline_step_tool_id ON pipeline_step (tool_id);
CREATE INDEX idx_run_organization_id ON run (organization_id);
CREATE INDEX idx_run_pipeline_id ON run (pipeline_id);
CREATE INDEX idx_run_tool_id ON run (tool_id);
CREATE INDEX idx_artifact_organization_id ON artifact (organization_id);
CREATE INDEX idx_artifact_producer_id ON artifact (producer_id);
CREATE UNIQUE INDEX idx_label_name_organization ON label (name, organization_id);
CREATE UNIQUE INDEX idx_organization_slug ON organization(slug);

-- Magic link token indexes
CREATE INDEX idx_magic_link_tokens_token_hash ON magic_link_tokens(token_hash);
CREATE INDEX idx_magic_link_tokens_code ON magic_link_tokens(code) WHERE code IS NOT NULL;
CREATE INDEX idx_magic_link_tokens_identifier ON magic_link_tokens(identifier);
CREATE INDEX idx_magic_link_tokens_expires_at ON magic_link_tokens(expires_at) WHERE used_at IS NULL;

-- Artifact embedding index for vector similarity search
CREATE INDEX ON artifact USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

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

-- Drop function
DROP FUNCTION IF EXISTS cleanup_expired_magic_links();

-- Drop junction tables first (they have foreign keys)
DROP TABLE IF EXISTS label_to_artifact CASCADE;
DROP TABLE IF EXISTS run_to_artifact CASCADE;
DROP TABLE IF EXISTS pipeline_step_to_dependency CASCADE;

-- Drop main tables in reverse dependency order
DROP TABLE IF EXISTS label CASCADE;
DROP TABLE IF EXISTS artifact CASCADE;
DROP TABLE IF EXISTS run CASCADE;
DROP TABLE IF EXISTS pipeline_step CASCADE;
DROP TABLE IF EXISTS tool CASCADE;
DROP TABLE IF EXISTS pipeline CASCADE;
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS magic_link_tokens CASCADE;
DROP TABLE IF EXISTS verification_token CASCADE;
DROP TABLE IF EXISTS invitation CASCADE;
DROP TABLE IF EXISTS member CASCADE;
DROP TABLE IF EXISTS session CASCADE;
DROP TABLE IF EXISTS account CASCADE;
DROP TABLE IF EXISTS organization CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;

-- Drop pgvector extension
DROP EXTENSION IF EXISTS vector;

-- +goose StatementEnd