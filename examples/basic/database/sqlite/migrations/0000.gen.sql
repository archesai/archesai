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
-- create "invitation" table
CREATE TABLE `invitation` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `email` text NOT NULL, `expires_at` text NOT NULL, `inviter_id` text NOT NULL, `organization_id` text NOT NULL, `role` text NOT NULL DEFAULT 'basic', `status` text NOT NULL DEFAULT 'pending', PRIMARY KEY (`id`), CONSTRAINT `invitation_inviter_id_fkey` FOREIGN KEY (`inviter_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `invitation_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_invitation_organization_id" to table: "invitation"
CREATE INDEX `idx_invitation_organization_id` ON `invitation` (`organization_id`);
-- create index "idx_invitation_inviter_id" to table: "invitation"
CREATE INDEX `idx_invitation_inviter_id` ON `invitation` (`inviter_id`);
-- create "member" table
CREATE TABLE `member` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `organization_id` text NOT NULL, `role` text NOT NULL DEFAULT 'basic', `user_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `member_organization_id_fkey` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `member_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_member_organization_id" to table: "member"
CREATE INDEX `idx_member_organization_id` ON `member` (`organization_id`);
-- create index "idx_member_user_id" to table: "member"
CREATE INDEX `idx_member_user_id` ON `member` (`user_id`);
-- create "organization" table
CREATE TABLE `organization` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `billing_email` text NULL, `credits` integer NOT NULL DEFAULT 0, `logo` text NULL, `name` text NOT NULL, `plan` text NOT NULL DEFAULT 'FREE', `slug` text NOT NULL, `stripe_customer_identifier` text NOT NULL, PRIMARY KEY (`id`));
-- create "session" table
CREATE TABLE `session` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `auth_method` text NULL, `auth_provider` text NULL, `expires_at` text NOT NULL, `ip_address` text NULL, `organization_id` text NULL, `token` text NOT NULL, `user_agent` text NULL, `user_id` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `session_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_session_user_id" to table: "session"
CREATE INDEX `idx_session_user_id` ON `session` (`user_id`);
-- create index "idx_session_token" to table: "session"
CREATE INDEX `idx_session_token` ON `session` (`token`);
-- create "todo" table
CREATE TABLE `todo` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `completed` integer NOT NULL, `title` text NOT NULL, PRIMARY KEY (`id`));
-- create "user" table
CREATE TABLE `user` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `email` text NOT NULL, `email_verified` integer NOT NULL DEFAULT 0, `image` text NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
