-- SQLite version of initial schema
-- Note: SQLite doesn't support ENUMs, so we use CHECK constraints instead

-- Create user table
CREATE TABLE "user" (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email TEXT NOT NULL UNIQUE,
    email_verified INTEGER NOT NULL DEFAULT 0, -- SQLite uses 0/1 for boolean
    image TEXT,
    name TEXT NOT NULL
);

-- Create organization table
CREATE TABLE organization (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    billing_email TEXT,
    credits INTEGER NOT NULL DEFAULT 0,
    logo TEXT,
    metadata TEXT,
    name TEXT NOT NULL,
    plan TEXT NOT NULL DEFAULT 'FREE' CHECK(plan IN ('BASIC', 'FREE', 'PREMIUM', 'STANDARD', 'UNLIMITED')),
    stripe_customer_id TEXT UNIQUE
);

-- Create account table
CREATE TABLE account (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    access_token TEXT,
    access_token_expires_at DATETIME,
    account_id TEXT NOT NULL,
    id_token TEXT,
    password TEXT,
    provider_id TEXT NOT NULL,
    refresh_token TEXT,
    refresh_token_expires_at DATETIME,
    scope TEXT,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE
);

-- Create member table
CREATE TABLE member (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    organization_id TEXT NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member' CHECK(role IN ('admin', 'owner', 'member')),
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    UNIQUE(organization_id, user_id)
);

-- Create api_token table
CREATE TABLE api_token (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    hashed_token TEXT NOT NULL UNIQUE,
    last_used_at DATETIME,
    name TEXT NOT NULL,
    organization_id TEXT NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    scopes TEXT NOT NULL DEFAULT '[]', -- JSON array as text
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE
);

-- Create verification_token table
CREATE TABLE verification_token (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    identifier TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    UNIQUE(identifier, token)
);

-- Create invitation table
CREATE TABLE invitation (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    organization_id TEXT NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member' CHECK(role IN ('admin', 'owner', 'member')),
    token TEXT NOT NULL UNIQUE,
    UNIQUE(email, organization_id)
);

-- Intelligence feature tables

-- Create tool table
CREATE TABLE tool (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    category TEXT NOT NULL,
    config TEXT, -- JSON config
    description TEXT,
    metadata TEXT, -- JSON metadata
    name TEXT NOT NULL,
    organization_id TEXT NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    UNIQUE(name, organization_id)
);

-- Create pipeline table
CREATE TABLE pipeline (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    config TEXT, -- JSON config
    description TEXT,
    metadata TEXT, -- JSON metadata
    name TEXT NOT NULL,
    organization_id TEXT NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    UNIQUE(name, organization_id)
);

-- Create pipeline_step table
CREATE TABLE pipeline_step (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    config TEXT, -- JSON config
    metadata TEXT, -- JSON metadata
    pipeline_id TEXT NOT NULL REFERENCES pipeline(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    tool_id TEXT NOT NULL REFERENCES tool(id) ON DELETE CASCADE,
    UNIQUE(pipeline_id, position)
);

-- Create run table
CREATE TABLE run (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    credits_used INTEGER NOT NULL DEFAULT 0,
    error TEXT,
    metadata TEXT, -- JSON metadata
    organization_id TEXT NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    parameters TEXT, -- JSON parameters
    pipeline_id TEXT NOT NULL REFERENCES pipeline(id) ON DELETE CASCADE,
    started_at DATETIME,
    status TEXT NOT NULL DEFAULT 'QUEUED' CHECK(status IN ('COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED')),
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE
);

-- Create artifact table
CREATE TABLE artifact (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    content_type TEXT NOT NULL,
    data BLOB,
    metadata TEXT, -- JSON metadata
    name TEXT NOT NULL,
    organization_id TEXT NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    run_id TEXT NOT NULL REFERENCES run(id) ON DELETE CASCADE,
    size INTEGER NOT NULL,
    storage_path TEXT
);

-- Create label table
CREATE TABLE label (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    artifact_id TEXT NOT NULL REFERENCES artifact(id) ON DELETE CASCADE,
    confidence REAL NOT NULL,
    metadata TEXT, -- JSON metadata
    value TEXT NOT NULL
);

-- Create indexes for performance
CREATE INDEX idx_user_email ON "user"(email);
CREATE INDEX idx_account_user_id ON account(user_id);
CREATE INDEX idx_account_provider_account ON account(provider_id, account_id);
CREATE INDEX idx_member_organization_id ON member(organization_id);
CREATE INDEX idx_member_user_id ON member(user_id);
CREATE INDEX idx_api_token_organization_id ON api_token(organization_id);
CREATE INDEX idx_api_token_user_id ON api_token(user_id);
CREATE INDEX idx_invitation_organization_id ON invitation(organization_id);
CREATE INDEX idx_invitation_token ON invitation(token);
CREATE INDEX idx_tool_organization_id ON tool(organization_id);
CREATE INDEX idx_pipeline_organization_id ON pipeline(organization_id);
CREATE INDEX idx_pipeline_step_pipeline_id ON pipeline_step(pipeline_id);
CREATE INDEX idx_run_organization_id ON run(organization_id);
CREATE INDEX idx_run_pipeline_id ON run(pipeline_id);
CREATE INDEX idx_run_user_id ON run(user_id);
CREATE INDEX idx_run_status ON run(status);
CREATE INDEX idx_artifact_run_id ON artifact(run_id);
CREATE INDEX idx_artifact_organization_id ON artifact(organization_id);
CREATE INDEX idx_label_artifact_id ON label(artifact_id);

-- Create triggers for updated_at
CREATE TRIGGER update_user_updated_at AFTER UPDATE ON "user"
    FOR EACH ROW BEGIN
        UPDATE "user" SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_organization_updated_at AFTER UPDATE ON organization
    FOR EACH ROW BEGIN
        UPDATE organization SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_account_updated_at AFTER UPDATE ON account
    FOR EACH ROW BEGIN
        UPDATE account SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_member_updated_at AFTER UPDATE ON member
    FOR EACH ROW BEGIN
        UPDATE member SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_api_token_updated_at AFTER UPDATE ON api_token
    FOR EACH ROW BEGIN
        UPDATE api_token SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_verification_token_updated_at AFTER UPDATE ON verification_token
    FOR EACH ROW BEGIN
        UPDATE verification_token SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_invitation_updated_at AFTER UPDATE ON invitation
    FOR EACH ROW BEGIN
        UPDATE invitation SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_tool_updated_at AFTER UPDATE ON tool
    FOR EACH ROW BEGIN
        UPDATE tool SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_pipeline_updated_at AFTER UPDATE ON pipeline
    FOR EACH ROW BEGIN
        UPDATE pipeline SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_pipeline_step_updated_at AFTER UPDATE ON pipeline_step
    FOR EACH ROW BEGIN
        UPDATE pipeline_step SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_run_updated_at AFTER UPDATE ON run
    FOR EACH ROW BEGIN
        UPDATE run SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_artifact_updated_at AFTER UPDATE ON artifact
    FOR EACH ROW BEGIN
        UPDATE artifact SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_label_updated_at AFTER UPDATE ON label
    FOR EACH ROW BEGIN
        UPDATE label SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;