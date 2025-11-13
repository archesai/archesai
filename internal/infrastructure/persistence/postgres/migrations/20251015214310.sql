-- Generated at: 2025-10-15T21:43:10-04:00
-- create "executor" table
CREATE TABLE "public"."executor" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "cpu_shares" integer NOT NULL DEFAULT 512, "dependencies" text NULL, "description" text NOT NULL, "env" text NULL, "execute_code" text NOT NULL, "extra_files" text NULL, "is_active" boolean NOT NULL DEFAULT true, "language" text NOT NULL, "memory_mb" integer NOT NULL DEFAULT 256, "name" text NOT NULL, "organization_id" uuid NOT NULL, "schema_in" text NULL, "schema_out" text NULL, "timeout" integer NOT NULL DEFAULT 30, "version" integer NOT NULL DEFAULT 1, PRIMARY KEY ("id"), CONSTRAINT "executor_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT "executor_language_check" CHECK (language = ANY (ARRAY['nodejs'::text, 'python'::text, 'go'::text])));
-- create index "idx_executor_language" to table: "executor"
CREATE INDEX "idx_executor_language" ON "public"."executor" ("language");
-- create index "idx_executor_organization_id" to table: "executor"
CREATE INDEX "idx_executor_organization_id" ON "public"."executor" ("organization_id");
