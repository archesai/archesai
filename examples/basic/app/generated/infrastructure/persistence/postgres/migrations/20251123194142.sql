-- Generated at: 2025-11-23T19:41:42-05:00
-- create "user" table
CREATE TABLE "public"."user" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "email" text NOT NULL, "name" text NOT NULL, PRIMARY KEY ("id"));
