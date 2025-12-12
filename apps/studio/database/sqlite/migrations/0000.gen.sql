-- create "account" table
CREATE TABLE `account` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `access_token` text NULL, `access_token_expires_at` text NULL, `account_identifier` text NOT NULL, `id_token` text NULL, `provider` text NOT NULL, `refresh_token` text NULL, `refresh_token_expires_at` text NULL, `scope` text NULL, `user_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `account_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_account_user_id" to table: "account"
CREATE INDEX `idx_account_user_id` ON `account` (`user_id`);
-- create "api_key" table
CREATE TABLE `api_key` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `expires_at` text NULL, `key_hash` text NOT NULL, `last_used_at` text NULL, `name` text NULL, `organization_id` text NOT NULL, `prefix` text NULL, `rate_limit` integer NOT NULL DEFAULT 60, `scopes` text NOT NULL DEFAULT '[]', `user_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `api_key_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `api_key_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_api_key_organization_id" to table: "api_key"
CREATE INDEX `idx_api_key_organization_id` ON `api_key` (`organization_id`);
-- create index "idx_api_key_user_id" to table: "api_key"
CREATE INDEX `idx_api_key_user_id` ON `api_key` (`user_id`);
-- create index "idx_api_key_key_hash" to table: "api_key"
CREATE INDEX `idx_api_key_key_hash` ON `api_key` (`key_hash`);
-- create index "idx_api_key_last_used_at" to table: "api_key"
CREATE INDEX `idx_api_key_last_used_at` ON `api_key` (`last_used_at`);
-- create "artifact" table
CREATE TABLE `artifact` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `credits` integer NOT NULL DEFAULT 0, `description` text NULL, `mime_type` text NOT NULL DEFAULT 'application/octet-stream', `name` text NULL, `organization_id` text NOT NULL, `preview_image` text NULL, `producer_id` text NULL, `text` text NULL, `url` text NULL, PRIMARY KEY (`id`), CONSTRAINT `artifact_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `artifact_producer_id_fkey` FOREIGN KEY (`producer_id`) REFERENCES `run` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- create index "idx_artifact_organization_id" to table: "artifact"
CREATE INDEX `idx_artifact_organization_id` ON `artifact` (`organization_id`);
-- create index "idx_artifact_producer_id" to table: "artifact"
CREATE INDEX `idx_artifact_producer_id` ON `artifact` (`producer_id`);
-- create "executor" table
CREATE TABLE `executor` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `cpu_shares` integer NOT NULL DEFAULT 512, `dependencies` text NULL, `description` text NOT NULL, `env` text NULL, `execute_code` text NOT NULL, `extra_files` text NULL, `is_active` integer NOT NULL DEFAULT 1, `language` text NOT NULL, `memory_mb` integer NOT NULL DEFAULT 256, `name` text NOT NULL, `organization_id` text NOT NULL, `schema_in` text NULL, `schema_out` text NULL, `timeout` integer NOT NULL DEFAULT 30, `version` integer NOT NULL DEFAULT 1, PRIMARY KEY (`id`), CONSTRAINT `executor_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "idx_executor_organization_id" to table: "executor"
CREATE INDEX `idx_executor_organization_id` ON `executor` (`organization_id`);
-- create index "idx_executor_language" to table: "executor"
CREATE INDEX `idx_executor_language` ON `executor` (`language`);
-- create "invitation" table
CREATE TABLE `invitation` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `email` text NOT NULL, `expires_at` text NOT NULL, `inviter_id` text NOT NULL, `organization_id` text NOT NULL, `role` text NOT NULL DEFAULT 'basic', `status` text NOT NULL DEFAULT 'pending', PRIMARY KEY (`id`), CONSTRAINT `invitation_inviter_id_fkey` FOREIGN KEY (`inviter_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `invitation_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_invitation_organization_id" to table: "invitation"
CREATE INDEX `idx_invitation_organization_id` ON `invitation` (`organization_id`);
-- create index "idx_invitation_inviter_id" to table: "invitation"
CREATE INDEX `idx_invitation_inviter_id` ON `invitation` (`inviter_id`);
-- create "label" table
CREATE TABLE `label` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `name` text NOT NULL, `organization_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `label_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "idx_label_organization_id" to table: "label"
CREATE INDEX `idx_label_organization_id` ON `label` (`organization_id`);
-- create index "idx_label_name" to table: "label"
CREATE INDEX `idx_label_name` ON `label` (`name`);
-- create "member" table
CREATE TABLE `member` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `organization_id` text NOT NULL, `role` text NOT NULL DEFAULT 'basic', `user_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `member_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `member_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_member_organization_id" to table: "member"
CREATE INDEX `idx_member_organization_id` ON `member` (`organization_id`);
-- create index "idx_member_user_id" to table: "member"
CREATE INDEX `idx_member_user_id` ON `member` (`user_id`);
-- create "organization" table
CREATE TABLE `organization` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `billing_email` text NULL, `credits` integer NOT NULL DEFAULT 0, `logo` text NULL, `name` text NOT NULL, `plan` text NOT NULL DEFAULT 'FREE', `slug` text NOT NULL, `stripe_customer_identifier` text NOT NULL, PRIMARY KEY (`id`));
-- create "pipeline" table
CREATE TABLE `pipeline` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `description` text NULL, `name` text NULL, `organization_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `pipeline_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "idx_pipeline_organization_id" to table: "pipeline"
CREATE INDEX `idx_pipeline_organization_id` ON `pipeline` (`organization_id`);
-- create "pipeline_step" table
CREATE TABLE `pipeline_step` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `pipeline_id` text NOT NULL, `tool_id` text NOT NULL, PRIMARY KEY (`id`));
-- create "run" table
CREATE TABLE `run` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `completed_at` text NULL, `error` text NULL, `organization_id` text NOT NULL, `pipeline_id` text NOT NULL, `progress` integer NOT NULL DEFAULT 0, `started_at` text NULL, `status` text NOT NULL DEFAULT 'QUEUED', `tool_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `run_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `run_pipeline_id_fkey` FOREIGN KEY (`pipeline_id`) REFERENCES `pipeline` (`id`) ON UPDATE CASCADE ON DELETE SET NULL, CONSTRAINT `run_tool_id_fkey` FOREIGN KEY (`tool_id`) REFERENCES `tool` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- create index "idx_run_pipeline_id" to table: "run"
CREATE INDEX `idx_run_pipeline_id` ON `run` (`pipeline_id`);
-- create index "idx_run_organization_id" to table: "run"
CREATE INDEX `idx_run_organization_id` ON `run` (`organization_id`);
-- create index "idx_run_tool_id" to table: "run"
CREATE INDEX `idx_run_tool_id` ON `run` (`tool_id`);
-- create "session" table
CREATE TABLE `session` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `auth_method` text NULL, `auth_provider` text NULL, `expires_at` text NOT NULL, `ip_address` text NULL, `organization_id` text NULL, `token` text NOT NULL, `user_agent` text NULL, `user_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `session_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_session_user_id" to table: "session"
CREATE INDEX `idx_session_user_id` ON `session` (`user_id`);
-- create index "idx_session_token" to table: "session"
CREATE INDEX `idx_session_token` ON `session` (`token`);
-- create "tool" table
CREATE TABLE `tool` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `description` text NOT NULL, `input_mime_type` text NOT NULL DEFAULT 'application/octet-stream', `name` text NOT NULL, `organization_id` text NOT NULL, `output_mime_type` text NOT NULL DEFAULT 'application/octet-stream', PRIMARY KEY (`id`), CONSTRAINT `tool_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "idx_tool_organization_id" to table: "tool"
CREATE INDEX `idx_tool_organization_id` ON `tool` (`organization_id`);
-- create "user" table
CREATE TABLE `user` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `email` text NOT NULL, `email_verified` integer NOT NULL DEFAULT 0, `image` text NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
