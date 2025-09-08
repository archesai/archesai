-- +goose Up
-- +goose StatementBegin

-- PostgreSQL schema for sqlc
-- This is the source of truth for the database schema

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
    metadata TEXT,
    name TEXT NOT NULL,
    plan TEXT NOT NULL DEFAULT 'FREE' CHECK (plan IN ('BASIC', 'FREE', 'PREMIUM', 'STANDARD', 'UNLIMITED')),
    stripe_customer_id TEXT UNIQUE
);

-- Create account table
CREATE TABLE account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    access_token TEXT,
    access_token_expires_at TIMESTAMPTZ,
    account_id TEXT NOT NULL,
    id_token TEXT,
    password TEXT,
    provider_id TEXT NOT NULL,
    refresh_token TEXT,
    refresh_token_expires_at TIMESTAMPTZ,
    scope TEXT,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE
);

-- Create session table
CREATE TABLE session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active_organization_id UUID,
    expires_at TIMESTAMPTZ NOT NULL,
    ip_address TEXT,
    token TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE
);

-- Create member table
CREATE TABLE member (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'owner', 'member')),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE
);

-- Create invitation table
CREATE TABLE invitation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    inviter_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
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

-- Create api_token table
CREATE TABLE api_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    enabled BOOLEAN NOT NULL,
    expires_at TIMESTAMPTZ,
    key TEXT NOT NULL,
    last_refill TIMESTAMPTZ,
    last_request TIMESTAMPTZ,
    metadata TEXT,
    name TEXT,
    permissions TEXT,
    prefix TEXT,
    rate_limit_enabled BOOLEAN NOT NULL,
    rate_limit_max INTEGER,
    rate_limit_time_window INTEGER,
    refill_amount INTEGER,
    refill_interval INTEGER,
    remaining INTEGER,
    request_count INTEGER NOT NULL DEFAULT 0,
    start TEXT,
    user_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create pipeline table
CREATE TABLE pipeline (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    name TEXT,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create tool table
CREATE TABLE tool (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT NOT NULL,
    input_mime_type TEXT NOT NULL DEFAULT 'application/octet-stream',
    name TEXT NOT NULL,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE ON UPDATE CASCADE,
    output_mime_type TEXT NOT NULL DEFAULT 'application/octet-stream'
);

-- Create pipeline_step table
CREATE TABLE pipeline_step (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    pipeline_id UUID NOT NULL REFERENCES pipeline(id) ON DELETE CASCADE ON UPDATE CASCADE,
    tool_id UUID NOT NULL REFERENCES tool(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create pipeline_step_to_dependency table (junction table)
CREATE TABLE pipeline_step_to_dependency (
    pipeline_step_id UUID NOT NULL REFERENCES pipeline_step(id) ON DELETE CASCADE,
    prerequisite_id UUID NOT NULL REFERENCES pipeline_step(id) ON DELETE CASCADE,
    PRIMARY KEY (pipeline_step_id, prerequisite_id)
);

-- Create run table
CREATE TABLE run (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    error TEXT,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE ON UPDATE CASCADE,
    pipeline_id UUID REFERENCES pipeline(id) ON DELETE SET NULL ON UPDATE CASCADE,
    progress DOUBLE PRECISION NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'QUEUED' CHECK (status IN ('COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED')),
    tool_id UUID REFERENCES tool(id) ON DELETE SET NULL ON UPDATE CASCADE
);

-- Create artifact table
CREATE TABLE artifact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    credits INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    mime_type TEXT NOT NULL DEFAULT 'application/octet-stream',
    name TEXT,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE ON UPDATE CASCADE,
    preview_image TEXT,
    producer_id UUID REFERENCES run(id) ON DELETE SET NULL ON UPDATE CASCADE,
    text TEXT,
    url TEXT
);

-- Create run_to_artifact table (junction table)
CREATE TABLE run_to_artifact (
    run_id UUID NOT NULL REFERENCES run(id) ON DELETE CASCADE,
    artifact_id UUID NOT NULL REFERENCES artifact(id) ON DELETE CASCADE,
    PRIMARY KEY (run_id, artifact_id)
);

-- Create label table
CREATE TABLE label (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create label_to_artifact table (junction table)
CREATE TABLE label_to_artifact (
    label_id UUID NOT NULL REFERENCES label(id) ON DELETE CASCADE,
    artifact_id UUID NOT NULL REFERENCES artifact(id) ON DELETE CASCADE,
    PRIMARY KEY (label_id, artifact_id)
);

-- Create indexes
CREATE INDEX idx_account_user_id ON account(user_id);
CREATE INDEX idx_session_user_id ON session(user_id);
CREATE INDEX idx_member_organization_id ON member(organization_id);
CREATE INDEX idx_member_user_id ON member(user_id);
CREATE INDEX idx_invitation_organization_id ON invitation(organization_id);
CREATE INDEX idx_invitation_inviter_id ON invitation(inviter_id);
CREATE INDEX idx_api_token_user_id ON api_token(user_id);
CREATE INDEX idx_pipeline_organization_id ON pipeline(organization_id);
CREATE INDEX idx_tool_organization_id ON tool(organization_id);
CREATE INDEX idx_pipeline_step_pipeline_id ON pipeline_step(pipeline_id);
CREATE INDEX idx_pipeline_step_tool_id ON pipeline_step(tool_id);
CREATE INDEX idx_run_organization_id ON run(organization_id);
CREATE INDEX idx_run_pipeline_id ON run(pipeline_id);
CREATE INDEX idx_run_tool_id ON run(tool_id);
CREATE INDEX idx_artifact_organization_id ON artifact(organization_id);
CREATE INDEX idx_artifact_producer_id ON artifact(producer_id);
CREATE UNIQUE INDEX idx_label_name_organization ON label(name, organization_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS artifact_label CASCADE;
DROP TABLE IF EXISTS label CASCADE;
DROP TABLE IF EXISTS artifact CASCADE;
DROP TABLE IF EXISTS run CASCADE;
DROP TABLE IF EXISTS pipeline_step CASCADE;
DROP TABLE IF EXISTS tool CASCADE;
DROP TABLE IF EXISTS pipeline CASCADE;
DROP TABLE IF EXISTS invitation CASCADE;
DROP TABLE IF EXISTS member CASCADE;
DROP TABLE IF EXISTS organization CASCADE;
DROP TABLE IF EXISTS session CASCADE;
DROP TABLE IF EXISTS account CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
-- +goose StatementEnd