-- Generated at: 2025-12-04T02:13:19-05:00
-- create "todo" table
CREATE TABLE `todo` (`id` text NOT NULL DEFAULT (lower(hex(randomblob(16)))), `created_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `updated_at` text NOT NULL DEFAULT (CURRENT_TIMESTAMP), `completed` integer NOT NULL, `title` text NOT NULL, PRIMARY KEY (`id`));
