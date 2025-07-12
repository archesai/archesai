CREATE TYPE "public"."authType" AS ENUM('email', 'oauth', 'oidc', 'webauthn');--> statement-breakpoint
CREATE TYPE "public"."planType" AS ENUM('BASIC', 'FREE', 'PREMIUM', 'STANDARD', 'UNLIMITED');--> statement-breakpoint
CREATE TYPE "public"."providerType" AS ENUM('API_KEY', 'FIREBASE', 'LOCAL', 'TWITTER');--> statement-breakpoint
CREATE TYPE "public"."role" AS ENUM('ADMIN', 'USER');--> statement-breakpoint
CREATE TYPE "public"."runStatus" AS ENUM('COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED');--> statement-breakpoint
CREATE TYPE "public"."RunType" AS ENUM('PIPELINE_RUN', 'TOOL_RUN');--> statement-breakpoint
CREATE TYPE "public"."toolIO" AS ENUM('AUDIO', 'IMAGE', 'TEXT', 'VIDEO');--> statement-breakpoint
CREATE TYPE "public"."verificationTokenType" AS ENUM('EMAIL_CHANGE', 'EMAIL_VERIFICATION', 'PASSWORD_RESET');--> statement-breakpoint
CREATE TABLE "accounts" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"access_token" text,
	"access_token_expires_at" date,
	"account_id" text NOT NULL,
	"id_token" text,
	"password" text,
	"provider_id" text NOT NULL,
	"refresh_token" text,
	"refresh_token_expires_at" date,
	"scope" text,
	"user_id" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "api-tokens" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"key" text NOT NULL,
	"organization_id" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL
);
--> statement-breakpoint
CREATE TABLE "artifacts" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"credits" integer DEFAULT 0 NOT NULL,
	"description" text,
	"embedding" vector(1536),
	"mime_type" text,
	"organization_id" text NOT NULL,
	"parent_id" text,
	"preview_image" text,
	"producer_id" text,
	"text" text,
	"url" text
);
--> statement-breakpoint
CREATE TABLE "_parentToChild" (
	"child_id" text NOT NULL,
	"parent_id" text NOT NULL,
	CONSTRAINT "_parentToChild_parent_id_child_id_pk" PRIMARY KEY("parent_id","child_id")
);
--> statement-breakpoint
CREATE TABLE "invitations" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"accepted" boolean DEFAULT false,
	"email" text NOT NULL,
	"organization_id" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL
);
--> statement-breakpoint
CREATE TABLE "labels" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"organization_id" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "labelToArtifact" (
	"artifact_id" text NOT NULL,
	"label_id" text NOT NULL,
	CONSTRAINT "labelToArtifact_label_id_artifact_id_pk" PRIMARY KEY("label_id","artifact_id")
);
--> statement-breakpoint
CREATE TABLE "members" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"invitationId" text,
	"organization_id" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL,
	"userId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "organizations" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"billing_email" text NOT NULL,
	"credits" integer DEFAULT 0 NOT NULL,
	"orgname" text NOT NULL,
	"plan" "planType" DEFAULT 'FREE' NOT NULL,
	"stripe_customer_id" text,
	CONSTRAINT "organizations_orgname_unique" UNIQUE("orgname"),
	CONSTRAINT "organizations_stripeCustomerId_unique" UNIQUE("stripe_customer_id")
);
--> statement-breakpoint
CREATE TABLE "pipelines" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"description" text,
	"organization_id" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "pipeline-steps" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"pipeline_id" text NOT NULL,
	"tool_id" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "pipelineStepToDependency" (
	"pipeline_step_id" text NOT NULL,
	"prerequisite_id" text NOT NULL,
	CONSTRAINT "pipelineStepToDependency_pipeline_step_id_prerequisite_id_pk" PRIMARY KEY("pipeline_step_id","prerequisite_id")
);
--> statement-breakpoint
CREATE TABLE "runs" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"completed_at" timestamp(3),
	"error" text,
	"organization_id" text NOT NULL,
	"pipeline_id" text,
	"progress" double precision DEFAULT 0 NOT NULL,
	"run_type" "RunType" NOT NULL,
	"started_at" timestamp(3),
	"status" "runStatus" DEFAULT 'QUEUED' NOT NULL,
	"tool_id" text
);
--> statement-breakpoint
CREATE TABLE "_runToContent" (
	"artifact_id" text NOT NULL,
	"run_id" text NOT NULL,
	CONSTRAINT "_runToContent_run_id_artifact_id_pk" PRIMARY KEY("run_id","artifact_id")
);
--> statement-breakpoint
CREATE TABLE "session" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"expires_at" date NOT NULL,
	"ip_address" text,
	"token" text NOT NULL,
	"user_agent" text,
	"userId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "tools" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"description" text NOT NULL,
	"input_type" "toolIO" NOT NULL,
	"organization_id" text NOT NULL,
	"output_type" "toolIO" NOT NULL,
	"tool_base" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "users" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"deactivated" boolean DEFAULT false NOT NULL,
	"email" text NOT NULL,
	"email_verified" boolean DEFAULT false NOT NULL,
	"image" text,
	"name" text NOT NULL,
	CONSTRAINT "users_email_unique" UNIQUE("email")
);
--> statement-breakpoint
CREATE TABLE "verification-tokens" (
	"created_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updated_at" date DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"expires_at" date NOT NULL,
	"identifier" text NOT NULL,
	"value" text NOT NULL
);
--> statement-breakpoint
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "api-tokens" ADD CONSTRAINT "api-tokens_organization_id_organizations_id_fk" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_organization_id_organizations_id_fk" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_parent_id_artifacts_id_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."artifacts"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_producer_id_runs_id_fk" FOREIGN KEY ("producer_id") REFERENCES "public"."runs"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_parentToChild" ADD CONSTRAINT "_parentToChild_child_id_artifacts_id_fk" FOREIGN KEY ("child_id") REFERENCES "public"."artifacts"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_parentToChild" ADD CONSTRAINT "_parentToChild_parent_id_artifacts_id_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."artifacts"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "invitations" ADD CONSTRAINT "invitations_organization_id_organizations_id_fk" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "labels" ADD CONSTRAINT "labels_organization_id_organizations_id_fk" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "labelToArtifact" ADD CONSTRAINT "labelToArtifact_artifact_id_artifacts_id_fk" FOREIGN KEY ("artifact_id") REFERENCES "public"."artifacts"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "labelToArtifact" ADD CONSTRAINT "labelToArtifact_label_id_labels_id_fk" FOREIGN KEY ("label_id") REFERENCES "public"."labels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_invitationId_invitations_id_fk" FOREIGN KEY ("invitationId") REFERENCES "public"."invitations"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_organization_id_organizations_id_fk" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "pipelines" ADD CONSTRAINT "pipelines_organization_id_organizations_id_fk" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipeline-steps" ADD CONSTRAINT "pipeline-steps_pipeline_id_pipelines_id_fk" FOREIGN KEY ("pipeline_id") REFERENCES "public"."pipelines"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipeline-steps" ADD CONSTRAINT "pipeline-steps_tool_id_tools_id_fk" FOREIGN KEY ("tool_id") REFERENCES "public"."tools"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipelineStepToDependency" ADD CONSTRAINT "pipelineStepToDependency_pipeline_step_id_pipeline-steps_id_fk" FOREIGN KEY ("pipeline_step_id") REFERENCES "public"."pipeline-steps"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "pipelineStepToDependency" ADD CONSTRAINT "pipelineStepToDependency_prerequisite_id_pipeline-steps_id_fk" FOREIGN KEY ("prerequisite_id") REFERENCES "public"."pipeline-steps"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_organization_id_organizations_id_fk" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_pipeline_id_pipelines_id_fk" FOREIGN KEY ("pipeline_id") REFERENCES "public"."pipelines"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_tool_id_tools_id_fk" FOREIGN KEY ("tool_id") REFERENCES "public"."tools"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_runToContent" ADD CONSTRAINT "_runToContent_artifact_id_artifacts_id_fk" FOREIGN KEY ("artifact_id") REFERENCES "public"."artifacts"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_runToContent" ADD CONSTRAINT "_runToContent_run_id_runs_id_fk" FOREIGN KEY ("run_id") REFERENCES "public"."runs"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "session" ADD CONSTRAINT "session_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "tools" ADD CONSTRAINT "tools_organization_id_organizations_id_fk" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
CREATE UNIQUE INDEX "labels_name_organization_id_index" ON "labels" USING btree ("name","organization_id");--> statement-breakpoint
CREATE UNIQUE INDEX "members_userId_organization_id_index" ON "members" USING btree ("userId","organization_id") WHERE "members"."userId" IS NOT NULL;