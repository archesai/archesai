-- create "pipeline_step" table
CREATE TABLE "public"."pipeline_step" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "pipeline_id" uuid NOT NULL, "tool_id" uuid NOT NULL, PRIMARY KEY ("id"));
-- create "user" table
CREATE TABLE "public"."user" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "email" text NOT NULL, "email_verified" boolean NOT NULL DEFAULT false, "image" text NULL, "name" text NOT NULL, PRIMARY KEY ("id"));
-- create "account" table
CREATE TABLE "public"."account" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "access_token" text NULL, "access_token_expires_at" timestamptz NULL, "account_identifier" text NOT NULL, "id_token" text NULL, "provider" text NOT NULL, "refresh_token" text NULL, "refresh_token_expires_at" timestamptz NULL, "scope" text NULL, "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "account_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "account_provider_check" CHECK (provider = ANY (ARRAY['local'::text, 'google'::text, 'github'::text, 'microsoft'::text, 'apple'::text])));
-- create index "idx_account_user_id" to table: "account"
CREATE INDEX "idx_account_user_id" ON "public"."account" ("user_id");
-- create "organization" table
CREATE TABLE "public"."organization" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "billing_email" text NULL, "credits" integer NOT NULL DEFAULT 0, "logo" text NULL, "name" text NOT NULL, "plan" text NOT NULL DEFAULT 'FREE', "slug" text NOT NULL, "stripe_customer_identifier" text NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "organization_plan_check" CHECK (plan = ANY (ARRAY['FREE'::text, 'BASIC'::text, 'STANDARD'::text, 'PREMIUM'::text, 'UNLIMITED'::text])));
-- create "api_key" table
CREATE TABLE "public"."api_key" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "expires_at" timestamptz NULL, "key_hash" text NOT NULL, "last_used_at" timestamptz NULL, "name" text NULL, "organization_id" uuid NOT NULL, "prefix" text NULL, "rate_limit" integer NOT NULL DEFAULT 60, "scopes" text[] NOT NULL DEFAULT '{}', "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "api_key_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "api_key_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "idx_api_key_key_hash" to table: "api_key"
CREATE INDEX "idx_api_key_key_hash" ON "public"."api_key" ("key_hash");
-- create index "idx_api_key_last_used_at" to table: "api_key"
CREATE INDEX "idx_api_key_last_used_at" ON "public"."api_key" ("last_used_at");
-- create index "idx_api_key_organization_id" to table: "api_key"
CREATE INDEX "idx_api_key_organization_id" ON "public"."api_key" ("organization_id");
-- create index "idx_api_key_user_id" to table: "api_key"
CREATE INDEX "idx_api_key_user_id" ON "public"."api_key" ("user_id");
-- create "pipeline" table
CREATE TABLE "public"."pipeline" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "description" text NULL, "name" text NULL, "organization_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "pipeline_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "idx_pipeline_organization_id" to table: "pipeline"
CREATE INDEX "idx_pipeline_organization_id" ON "public"."pipeline" ("organization_id");
-- create "tool" table
CREATE TABLE "public"."tool" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "description" text NOT NULL, "input_mime_type" text NOT NULL DEFAULT 'application/octet-stream', "name" text NOT NULL, "organization_id" uuid NOT NULL, "output_mime_type" text NOT NULL DEFAULT 'application/octet-stream', PRIMARY KEY ("id"), CONSTRAINT "tool_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "idx_tool_organization_id" to table: "tool"
CREATE INDEX "idx_tool_organization_id" ON "public"."tool" ("organization_id");
-- create "run" table
CREATE TABLE "public"."run" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "completed_at" timestamptz NULL, "error" text NULL, "organization_id" uuid NOT NULL, "pipeline_id" uuid NOT NULL, "progress" integer NOT NULL DEFAULT 0, "started_at" timestamptz NULL, "status" text NOT NULL DEFAULT 'QUEUED', "tool_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "run_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT "run_pipeline_id_fkey" FOREIGN KEY ("pipeline_id") REFERENCES "public"."pipeline" ("id") ON UPDATE CASCADE ON DELETE SET NULL, CONSTRAINT "run_tool_id_fkey" FOREIGN KEY ("tool_id") REFERENCES "public"."tool" ("id") ON UPDATE CASCADE ON DELETE SET NULL, CONSTRAINT "run_status_check" CHECK (status = ANY (ARRAY['COMPLETED'::text, 'FAILED'::text, 'PROCESSING'::text, 'QUEUED'::text])));
-- create index "idx_run_organization_id" to table: "run"
CREATE INDEX "idx_run_organization_id" ON "public"."run" ("organization_id");
-- create index "idx_run_pipeline_id" to table: "run"
CREATE INDEX "idx_run_pipeline_id" ON "public"."run" ("pipeline_id");
-- create index "idx_run_tool_id" to table: "run"
CREATE INDEX "idx_run_tool_id" ON "public"."run" ("tool_id");
-- create "artifact" table
CREATE TABLE "public"."artifact" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "credits" integer NOT NULL DEFAULT 0, "description" text NULL, "mime_type" text NOT NULL DEFAULT 'application/octet-stream', "name" text NULL, "organization_id" uuid NOT NULL, "preview_image" text NULL, "producer_id" uuid NULL, "text" text NULL, "url" text NULL, PRIMARY KEY ("id"), CONSTRAINT "artifact_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT "artifact_producer_id_fkey" FOREIGN KEY ("producer_id") REFERENCES "public"."run" ("id") ON UPDATE CASCADE ON DELETE SET NULL);
-- create index "idx_artifact_organization_id" to table: "artifact"
CREATE INDEX "idx_artifact_organization_id" ON "public"."artifact" ("organization_id");
-- create index "idx_artifact_producer_id" to table: "artifact"
CREATE INDEX "idx_artifact_producer_id" ON "public"."artifact" ("producer_id");
-- create "executor" table
CREATE TABLE "public"."executor" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "cpu_shares" integer NOT NULL DEFAULT 512, "dependencies" text NULL, "description" text NOT NULL, "env" text NULL, "execute_code" text NOT NULL, "extra_files" text NULL, "is_active" boolean NOT NULL DEFAULT true, "language" text NOT NULL, "memory_mb" integer NOT NULL DEFAULT 256, "name" text NOT NULL, "organization_id" uuid NOT NULL, "schema_in" text NULL, "schema_out" text NULL, "timeout" integer NOT NULL DEFAULT 30, "version" integer NOT NULL DEFAULT 1, PRIMARY KEY ("id"), CONSTRAINT "executor_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT "executor_language_check" CHECK (language = ANY (ARRAY['nodejs'::text, 'python'::text, 'go'::text])));
-- create index "idx_executor_language" to table: "executor"
CREATE INDEX "idx_executor_language" ON "public"."executor" ("language");
-- create index "idx_executor_organization_id" to table: "executor"
CREATE INDEX "idx_executor_organization_id" ON "public"."executor" ("organization_id");
-- create "invitation" table
CREATE TABLE "public"."invitation" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "email" text NOT NULL, "expires_at" timestamptz NOT NULL, "inviter_id" uuid NOT NULL, "organization_id" uuid NOT NULL, "role" text NOT NULL DEFAULT 'basic', "status" text NOT NULL DEFAULT 'pending', PRIMARY KEY ("id"), CONSTRAINT "invitation_inviter_id_fkey" FOREIGN KEY ("inviter_id") REFERENCES "public"."user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "invitation_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "invitation_role_check" CHECK (role = ANY (ARRAY['admin'::text, 'owner'::text, 'basic'::text])), CONSTRAINT "invitation_status_check" CHECK (status = ANY (ARRAY['pending'::text, 'accepted'::text, 'declined'::text, 'expired'::text])));
-- create index "idx_invitation_inviter_id" to table: "invitation"
CREATE INDEX "idx_invitation_inviter_id" ON "public"."invitation" ("inviter_id");
-- create index "idx_invitation_organization_id" to table: "invitation"
CREATE INDEX "idx_invitation_organization_id" ON "public"."invitation" ("organization_id");
-- create "label" table
CREATE TABLE "public"."label" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "name" text NOT NULL, "organization_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "label_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "idx_label_name" to table: "label"
CREATE INDEX "idx_label_name" ON "public"."label" ("name");
-- create index "idx_label_organization_id" to table: "label"
CREATE INDEX "idx_label_organization_id" ON "public"."label" ("organization_id");
-- create "member" table
CREATE TABLE "public"."member" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "organization_id" uuid NOT NULL, "role" text NOT NULL DEFAULT 'basic', "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "member_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organization" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "member_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "member_role_check" CHECK (role = ANY (ARRAY['admin'::text, 'owner'::text, 'basic'::text])));
-- create index "idx_member_organization_id" to table: "member"
CREATE INDEX "idx_member_organization_id" ON "public"."member" ("organization_id");
-- create index "idx_member_user_id" to table: "member"
CREATE INDEX "idx_member_user_id" ON "public"."member" ("user_id");
-- create "session" table
CREATE TABLE "public"."session" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "auth_method" text NULL, "auth_provider" text NULL, "expires_at" timestamptz NOT NULL, "ip_address" text NULL, "organization_id" uuid NULL, "token" text NOT NULL, "user_agent" text NULL, "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "session_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "session_auth_provider_check" CHECK (auth_provider = ANY (ARRAY['local'::text, 'google'::text, 'github'::text, 'microsoft'::text, 'apple'::text])));
-- create index "idx_session_token" to table: "session"
CREATE INDEX "idx_session_token" ON "public"."session" ("token");
-- create index "idx_session_user_id" to table: "session"
CREATE INDEX "idx_session_user_id" ON "public"."session" ("user_id");
