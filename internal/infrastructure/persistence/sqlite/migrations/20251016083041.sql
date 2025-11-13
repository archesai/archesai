-- Generated at: 2025-10-16T08:30:41-04:00
-- drop index "idx_organization_slug" from table: "organization"
DROP INDEX `idx_organization_slug`;
-- create index "idx_session_token" to table: "session"
CREATE INDEX `idx_session_token` ON `session` (`token`);
