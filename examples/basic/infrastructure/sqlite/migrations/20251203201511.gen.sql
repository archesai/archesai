-- Generated at: 2025-12-03T20:15:11-05:00
-- create "todo" table
CREATE TABLE `todo` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `completed` integer NOT NULL, `title` text NOT NULL, PRIMARY KEY (`id`));
