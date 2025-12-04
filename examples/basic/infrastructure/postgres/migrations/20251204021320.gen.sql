-- Generated at: 2025-12-04T02:13:20-05:00
-- create "todo" table
CREATE TABLE "public"."todo" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "completed" boolean NOT NULL, "title" text NOT NULL, PRIMARY KEY ("id"));
