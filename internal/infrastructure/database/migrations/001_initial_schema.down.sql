-- Drop triggers
DROP TRIGGER IF EXISTS update_label_updated_at ON label;
DROP TRIGGER IF EXISTS update_artifact_updated_at ON artifact;
DROP TRIGGER IF EXISTS update_run_updated_at ON run;
DROP TRIGGER IF EXISTS update_pipeline_step_updated_at ON pipeline_step;
DROP TRIGGER IF EXISTS update_tool_updated_at ON tool;
DROP TRIGGER IF EXISTS update_pipeline_updated_at ON pipeline;
DROP TRIGGER IF EXISTS update_api_token_updated_at ON api_token;
DROP TRIGGER IF EXISTS update_verification_token_updated_at ON verification_token;
DROP TRIGGER IF EXISTS update_invitation_updated_at ON invitation;
DROP TRIGGER IF EXISTS update_member_updated_at ON member;
DROP TRIGGER IF EXISTS update_session_updated_at ON session;
DROP TRIGGER IF EXISTS update_account_updated_at ON account;
DROP TRIGGER IF EXISTS update_organization_updated_at ON organization;
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_label_slug_organization;
DROP INDEX IF EXISTS idx_artifact_slug_organization;
DROP INDEX IF EXISTS idx_run_slug_organization;
DROP INDEX IF EXISTS idx_pipeline_step_slug_pipeline;
DROP INDEX IF EXISTS idx_tool_slug_organization;
DROP INDEX IF EXISTS idx_pipeline_slug_organization;
DROP INDEX IF EXISTS idx_api_token_slug_organization;
DROP INDEX IF EXISTS idx_invitation_slug_organization;
DROP INDEX IF EXISTS idx_member_slug_organization;
DROP INDEX IF EXISTS idx_label_name_organization;
DROP INDEX IF EXISTS idx_artifact_producer_id;
DROP INDEX IF EXISTS idx_artifact_organization_id;
DROP INDEX IF EXISTS idx_run_tool_id;
DROP INDEX IF EXISTS idx_run_pipeline_id;
DROP INDEX IF EXISTS idx_run_organization_id;
DROP INDEX IF EXISTS idx_pipeline_step_tool_id;
DROP INDEX IF EXISTS idx_pipeline_step_pipeline_id;
DROP INDEX IF EXISTS idx_tool_organization_id;
DROP INDEX IF EXISTS idx_pipeline_organization_id;
DROP INDEX IF EXISTS idx_api_token_user_id;
DROP INDEX IF EXISTS idx_invitation_inviter_id;
DROP INDEX IF EXISTS idx_invitation_organization_id;
DROP INDEX IF EXISTS idx_member_user_id;
DROP INDEX IF EXISTS idx_member_organization_id;
DROP INDEX IF EXISTS idx_session_user_id;
DROP INDEX IF EXISTS idx_account_user_id;

-- Drop junction tables
DROP TABLE IF EXISTS label_to_artifact;
DROP TABLE IF EXISTS run_to_artifact;
DROP TABLE IF EXISTS pipeline_step_to_dependency;

-- Drop tables (in dependency order)
DROP TABLE IF EXISTS label;
DROP TABLE IF EXISTS artifact;
DROP TABLE IF EXISTS run;
DROP TABLE IF EXISTS pipeline_step;
DROP TABLE IF EXISTS tool;
DROP TABLE IF EXISTS pipeline;
DROP TABLE IF EXISTS api_token;
DROP TABLE IF EXISTS verification_token;
DROP TABLE IF EXISTS invitation;
DROP TABLE IF EXISTS member;
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS account;
DROP TABLE IF EXISTS organization;
DROP TABLE IF EXISTS "user";

-- Drop enums
DROP TYPE IF EXISTS run_status;
DROP TYPE IF EXISTS plan_type;
DROP TYPE IF EXISTS role;

-- Drop extensions (optional, you may want to keep these)
-- DROP EXTENSION IF EXISTS "vector";
-- DROP EXTENSION IF EXISTS "uuid-ossp";