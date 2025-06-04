CREATE TYPE "public"."authType" AS ENUM('email', 'oauth', 'oidc', 'webauthn');--> statement-breakpoint
CREATE TYPE "public"."planType" AS ENUM('BASIC', 'FREE', 'PREMIUM', 'STANDARD', 'UNLIMITED');--> statement-breakpoint
CREATE TYPE "public"."providerType" AS ENUM('API_KEY', 'FIREBASE', 'LOCAL', 'TWITTER');--> statement-breakpoint
CREATE TYPE "public"."role" AS ENUM('ADMIN', 'USER');--> statement-breakpoint
CREATE TYPE "public"."runStatus" AS ENUM('COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED');--> statement-breakpoint
CREATE TYPE "public"."RunType" AS ENUM('PIPELINE_RUN', 'TOOL_RUN');--> statement-breakpoint
CREATE TYPE "public"."toolIO" AS ENUM('AUDIO', 'IMAGE', 'TEXT', 'VIDEO');--> statement-breakpoint
CREATE TYPE "public"."verificationTokenType" AS ENUM('EMAIL_CHANGE', 'EMAIL_VERIFICATION', 'PASSWORD_RESET');--> statement-breakpoint
CREATE TABLE "accounts" (
	"access_token" text,
	"expires_at" integer,
	"hashed_password" text,
	"id" text DEFAULT gen_random_uuid() NOT NULL,
	"id_token" text,
	"provider" "providerType" NOT NULL,
	"providerAccountId" text NOT NULL,
	"refresh_token" text,
	"scope" text,
	"session_state" text,
	"token_type" text,
	"type" "authType" NOT NULL,
	"userId" text NOT NULL,
	CONSTRAINT "accounts_provider_providerAccountId_pk" PRIMARY KEY("provider","providerAccountId"),
	CONSTRAINT "accounts_id_unique" UNIQUE("id")
);
--> statement-breakpoint
CREATE TABLE "api-tokens" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"key" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL
);
--> statement-breakpoint
CREATE TABLE "authenticator" (
	"counter" integer NOT NULL,
	"credentialBackedUp" boolean NOT NULL,
	"credentialDeviceType" text NOT NULL,
	"credentialID" text NOT NULL,
	"credentialPublicKey" text NOT NULL,
	"providerAccountId" text NOT NULL,
	"transports" text,
	"userId" text NOT NULL,
	CONSTRAINT "authenticator_userId_credentialID_pk" PRIMARY KEY("userId","credentialID"),
	CONSTRAINT "authenticator_credentialID_unique" UNIQUE("credentialID")
);
--> statement-breakpoint
CREATE TABLE "contents" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"credits" integer DEFAULT 0 NOT NULL,
	"description" text,
	"embedding" vector(1536),
	"mime_type" text,
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
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"accepted" boolean DEFAULT false,
	"email" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL
);
--> statement-breakpoint
CREATE TABLE "labels" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL
);
--> statement-breakpoint
CREATE TABLE "_labelsToContent" (
	"content_id" text NOT NULL,
	"label_id" text NOT NULL,
	CONSTRAINT "_labelsToContent_label_id_content_id_pk" PRIMARY KEY("label_id","content_id")
);
--> statement-breakpoint
CREATE TABLE "members" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"invitationId" text,
	"role" "role" DEFAULT 'USER' NOT NULL,
	"userId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "organizations" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"billing_email" text NOT NULL,
	"credits" integer DEFAULT 0 NOT NULL,
	"plan" "planType" DEFAULT 'FREE' NOT NULL,
	"stripe_customer_id" text,
	CONSTRAINT "organizations_orgname_unique" UNIQUE("orgname"),
	CONSTRAINT "organizations_stripeCustomerId_unique" UNIQUE("stripe_customer_id")
);
--> statement-breakpoint
CREATE TABLE "pipeline-steps" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"pipeline_id" text NOT NULL,
	"tool_id" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "pipelines" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"description" text
);
--> statement-breakpoint
CREATE TABLE "_pipelineStepDependencies" (
	"pipeline_step_id" text NOT NULL,
	"prerequisite_step_id" text NOT NULL,
	CONSTRAINT "_pipelineStepDependencies_pipeline_step_id_prerequisite_step_id_pk" PRIMARY KEY("pipeline_step_id","prerequisite_step_id")
);
--> statement-breakpoint
CREATE TABLE "runs" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"completed_at" timestamp(3),
	"error" text,
	"pipeline_id" text,
	"progress" double precision DEFAULT 0 NOT NULL,
	"run_type" "RunType" NOT NULL,
	"started_at" timestamp(3),
	"status" "runStatus" DEFAULT 'QUEUED' NOT NULL,
	"tool_id" text
);
--> statement-breakpoint
CREATE TABLE "_runToContent" (
	"content_id" text NOT NULL,
	"run_id" text NOT NULL,
	CONSTRAINT "_runToContent_run_id_content_id_pk" PRIMARY KEY("run_id","content_id")
);
--> statement-breakpoint
CREATE TABLE "session" (
	"expires" timestamp NOT NULL,
	"sessionToken" text PRIMARY KEY NOT NULL,
	"userId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "tools" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"description" text NOT NULL,
	"input_type" "toolIO" NOT NULL,
	"output_type" "toolIO" NOT NULL,
	"tool_base" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "users" (
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text,
	"orgname" text NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"deactivated" boolean DEFAULT false NOT NULL,
	"email" text,
	"emailVerified" timestamp,
	"image" text,
	CONSTRAINT "users_email_unique" UNIQUE("email")
);
--> statement-breakpoint
CREATE TABLE "verification-tokens" (
	"expires" timestamp NOT NULL,
	"id" text DEFAULT gen_random_uuid() NOT NULL,
	"identifier" text NOT NULL,
	"newEmail" text,
	"token" text NOT NULL,
	"type" "verificationTokenType" NOT NULL,
	CONSTRAINT "verification-tokens_id_unique" UNIQUE("id")
);
--> statement-breakpoint
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "api-tokens" ADD CONSTRAINT "api-tokens_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "authenticator" ADD CONSTRAINT "authenticator_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "contents" ADD CONSTRAINT "contents_parent_id_contents_id_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."contents"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "contents" ADD CONSTRAINT "contents_producer_id_runs_id_fk" FOREIGN KEY ("producer_id") REFERENCES "public"."runs"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "contents" ADD CONSTRAINT "contents_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_parentToChild" ADD CONSTRAINT "_parentToChild_child_id_contents_id_fk" FOREIGN KEY ("child_id") REFERENCES "public"."contents"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_parentToChild" ADD CONSTRAINT "_parentToChild_parent_id_contents_id_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."contents"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "invitations" ADD CONSTRAINT "invitations_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "labels" ADD CONSTRAINT "labels_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_labelsToContent" ADD CONSTRAINT "_labelsToContent_content_id_contents_id_fk" FOREIGN KEY ("content_id") REFERENCES "public"."contents"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_labelsToContent" ADD CONSTRAINT "_labelsToContent_label_id_labels_id_fk" FOREIGN KEY ("label_id") REFERENCES "public"."labels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_invitationId_invitations_id_fk" FOREIGN KEY ("invitationId") REFERENCES "public"."invitations"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipeline-steps" ADD CONSTRAINT "pipeline-steps_pipeline_id_pipelines_id_fk" FOREIGN KEY ("pipeline_id") REFERENCES "public"."pipelines"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipeline-steps" ADD CONSTRAINT "pipeline-steps_tool_id_tools_id_fk" FOREIGN KEY ("tool_id") REFERENCES "public"."tools"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipelines" ADD CONSTRAINT "pipelines_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_pipelineStepDependencies" ADD CONSTRAINT "_pipelineStepDependencies_pipeline_step_id_pipeline-steps_id_fk" FOREIGN KEY ("pipeline_step_id") REFERENCES "public"."pipeline-steps"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_pipelineStepDependencies" ADD CONSTRAINT "_pipelineStepDependencies_prerequisite_step_id_pipeline-steps_id_fk" FOREIGN KEY ("prerequisite_step_id") REFERENCES "public"."pipeline-steps"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_pipeline_id_pipelines_id_fk" FOREIGN KEY ("pipeline_id") REFERENCES "public"."pipelines"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_tool_id_tools_id_fk" FOREIGN KEY ("tool_id") REFERENCES "public"."tools"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_runToContent" ADD CONSTRAINT "_runToContent_content_id_contents_id_fk" FOREIGN KEY ("content_id") REFERENCES "public"."contents"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_runToContent" ADD CONSTRAINT "_runToContent_run_id_runs_id_fk" FOREIGN KEY ("run_id") REFERENCES "public"."runs"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "session" ADD CONSTRAINT "session_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "tools" ADD CONSTRAINT "tools_orgname_organizations_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
CREATE UNIQUE INDEX "labels_name_orgname_index" ON "labels" USING btree ("name","orgname");--> statement-breakpoint
CREATE UNIQUE INDEX "members_userId_orgname_index" ON "members" USING btree ("userId","orgname") WHERE "members"."userId" IS NOT NULL;--> statement-breakpoint
CREATE UNIQUE INDEX "pipeline-steps_name_pipeline_id_index" ON "pipeline-steps" USING btree ("name","pipeline_id");