-- Generated at: 2025-11-23T19:41:42-05:00
-- create "user" table
CREATE TABLE `user` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `email` text NOT NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
