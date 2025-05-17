CREATE TYPE "public"."authType" AS ENUM('email', 'oidc', 'webauthn', 'oauth');--> statement-breakpoint
CREATE TYPE "public"."planType" AS ENUM('FREE', 'BASIC', 'STANDARD', 'PREMIUM', 'UNLIMITED');--> statement-breakpoint
CREATE TYPE "public"."providerType" AS ENUM('firebase', 'credentials', 'twitter');--> statement-breakpoint
CREATE TYPE "public"."role" AS ENUM('USER', 'ADMIN');--> statement-breakpoint
CREATE TYPE "public"."StatusTypeEnum" AS ENUM('QUEUED', 'PROCESSING', 'COMPLETED', 'FAILED');--> statement-breakpoint
CREATE TYPE "public"."RunType" AS ENUM('PIPELINE_RUN', 'TOOL_RUN');--> statement-breakpoint
CREATE TYPE "public"."toolIO" AS ENUM('TEXT', 'IMAGE', 'VIDEO', 'AUDIO');--> statement-breakpoint
CREATE TYPE "public"."verificationTokenType" AS ENUM('EMAIL_VERIFICATION', 'PASSWORD_RESET', 'EMAIL_CHANGE');--> statement-breakpoint
CREATE TABLE "account" (
	"userId" text NOT NULL,
	"type" "authType" NOT NULL,
	"provider" "providerType" NOT NULL,
	"providerAccountId" text NOT NULL,
	"refresh_token" text,
	"access_token" text,
	"expires_at" integer,
	"token_type" text,
	"scope" text,
	"id_token" text,
	"session_state" text,
	"hashed_password" text,
	"id" text DEFAULT gen_random_uuid() NOT NULL,
	CONSTRAINT "account_provider_providerAccountId_pk" PRIMARY KEY("provider","providerAccountId"),
	CONSTRAINT "account_id_unique" UNIQUE("id")
);
--> statement-breakpoint
CREATE TABLE "authenticator" (
	"credentialID" text NOT NULL,
	"userId" text NOT NULL,
	"providerAccountId" text NOT NULL,
	"credentialPublicKey" text NOT NULL,
	"counter" integer NOT NULL,
	"credentialDeviceType" text NOT NULL,
	"credentialBackedUp" boolean NOT NULL,
	"transports" text,
	"id" text DEFAULT gen_random_uuid() NOT NULL,
	CONSTRAINT "authenticator_userId_credentialID_pk" PRIMARY KEY("userId","credentialID"),
	CONSTRAINT "authenticator_credentialID_unique" UNIQUE("credentialID"),
	CONSTRAINT "authenticator_id_unique" UNIQUE("id")
);
--> statement-breakpoint
CREATE TABLE "session" (
	"sessionToken" text PRIMARY KEY NOT NULL,
	"userId" text NOT NULL,
	"expires" timestamp NOT NULL,
	"id" text DEFAULT gen_random_uuid() NOT NULL,
	CONSTRAINT "session_id_unique" UNIQUE("id")
);
--> statement-breakpoint
CREATE TABLE "verificationToken" (
	"identifier" text NOT NULL,
	"token" text NOT NULL,
	"expires" timestamp NOT NULL,
	"type" "verificationTokenType" NOT NULL,
	"newEmail" text,
	"id" text DEFAULT gen_random_uuid() NOT NULL,
	CONSTRAINT "verificationToken_id_unique" UNIQUE("id")
);
--> statement-breakpoint
CREATE TABLE "apiToken" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL,
	"key" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "content" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"url" text,
	"text" text,
	"description" text,
	"embedding" vector(1536),
	"credits" integer DEFAULT 0 NOT NULL,
	"mime_type" text,
	"preview_image" text,
	"parent_id" text,
	"producer_id" text
);
--> statement-breakpoint
CREATE TABLE "label" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "_labelsToContent" (
	"label_id" text NOT NULL,
	"content_id" text NOT NULL,
	CONSTRAINT "_labelsToContent_label_id_content_id_pk" PRIMARY KEY("label_id","content_id")
);
--> statement-breakpoint
CREATE TABLE "member" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"invite_accepted" boolean DEFAULT false NOT NULL,
	"invite_email" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL,
	"userId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "organization" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"stripe_customer_id" text,
	"billing_email" text NOT NULL,
	"credits" integer DEFAULT 0 NOT NULL,
	"plan" "planType" DEFAULT 'FREE' NOT NULL,
	CONSTRAINT "organization_orgname_unique" UNIQUE("orgname"),
	CONSTRAINT "organization_stripeCustomerId_unique" UNIQUE("stripe_customer_id")
);
--> statement-breakpoint
CREATE TABLE "pipelineStep" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"pipeline_id" text NOT NULL,
	"tool_id" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "pipeline" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"description" text
);
--> statement-breakpoint
CREATE TABLE "_pipelineStepDependencies" (
	"pipeline_step_id" text NOT NULL,
	"prerequisite_step_id" text NOT NULL,
	CONSTRAINT "_pipelineStepDependencies_pipeline_step_id_prerequisite_step_id_pk" PRIMARY KEY("pipeline_step_id","prerequisite_step_id")
);
--> statement-breakpoint
CREATE TABLE "run" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"status" "StatusTypeEnum" DEFAULT 'QUEUED' NOT NULL,
	"started_at" timestamp(3),
	"completed_at" timestamp(3),
	"progress" double precision DEFAULT 0 NOT NULL,
	"error" text,
	"run_type" "RunType" NOT NULL,
	"pipeline_id" text,
	"tool_id" text,
	"pipeline_run_id" text NOT NULL,
	"pipeline_step_id" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "_runToContent" (
	"run_id" text NOT NULL,
	"content_id" text NOT NULL,
	CONSTRAINT "_runToContent_run_id_content_id_pk" PRIMARY KEY("run_id","content_id")
);
--> statement-breakpoint
CREATE TABLE "tool" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"description" text NOT NULL,
	"input_type" "toolIO" NOT NULL,
	"output_type" "toolIO" NOT NULL,
	"tool_base" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "user" (
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updated_at" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" text NOT NULL,
	"orgname" text NOT NULL,
	"email" text NOT NULL,
	"emailVerified" timestamp,
	"deactivated" boolean DEFAULT false NOT NULL,
	"image" text,
	CONSTRAINT "user_email_unique" UNIQUE("email")
);
--> statement-breakpoint
ALTER TABLE "account" ADD CONSTRAINT "account_userId_user_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."user"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "authenticator" ADD CONSTRAINT "authenticator_userId_user_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."user"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "session" ADD CONSTRAINT "session_userId_user_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."user"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "apiToken" ADD CONSTRAINT "apiToken_orgname_organization_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organization"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "content" ADD CONSTRAINT "content_parent_id_content_id_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."content"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "content" ADD CONSTRAINT "content_producer_id_run_id_fk" FOREIGN KEY ("producer_id") REFERENCES "public"."run"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "content" ADD CONSTRAINT "content_orgname_organization_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organization"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "label" ADD CONSTRAINT "label_orgname_organization_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organization"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_labelsToContent" ADD CONSTRAINT "_labelsToContent_label_id_label_id_fk" FOREIGN KEY ("label_id") REFERENCES "public"."label"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_labelsToContent" ADD CONSTRAINT "_labelsToContent_content_id_content_id_fk" FOREIGN KEY ("content_id") REFERENCES "public"."content"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "member" ADD CONSTRAINT "member_userId_user_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."user"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "member" ADD CONSTRAINT "member_orgname_organization_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organization"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipelineStep" ADD CONSTRAINT "pipelineStep_pipeline_id_pipeline_id_fk" FOREIGN KEY ("pipeline_id") REFERENCES "public"."pipeline"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipelineStep" ADD CONSTRAINT "pipelineStep_tool_id_tool_id_fk" FOREIGN KEY ("tool_id") REFERENCES "public"."tool"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipeline" ADD CONSTRAINT "pipeline_orgname_organization_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organization"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_pipelineStepDependencies" ADD CONSTRAINT "_pipelineStepDependencies_pipeline_step_id_pipelineStep_id_fk" FOREIGN KEY ("pipeline_step_id") REFERENCES "public"."pipelineStep"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_pipelineStepDependencies" ADD CONSTRAINT "_pipelineStepDependencies_prerequisite_step_id_pipelineStep_id_fk" FOREIGN KEY ("prerequisite_step_id") REFERENCES "public"."pipelineStep"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "run" ADD CONSTRAINT "run_pipeline_id_pipeline_id_fk" FOREIGN KEY ("pipeline_id") REFERENCES "public"."pipeline"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "run" ADD CONSTRAINT "run_tool_id_tool_id_fk" FOREIGN KEY ("tool_id") REFERENCES "public"."tool"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "run" ADD CONSTRAINT "run_orgname_organization_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organization"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "_runToContent" ADD CONSTRAINT "_runToContent_run_id_run_id_fk" FOREIGN KEY ("run_id") REFERENCES "public"."run"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "_runToContent" ADD CONSTRAINT "_runToContent_content_id_content_id_fk" FOREIGN KEY ("content_id") REFERENCES "public"."content"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "tool" ADD CONSTRAINT "tool_orgname_organization_id_fk" FOREIGN KEY ("orgname") REFERENCES "public"."organization"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
CREATE UNIQUE INDEX "label_name_orgname_index" ON "label" USING btree ("name","orgname");--> statement-breakpoint
CREATE UNIQUE INDEX "member_userId_orgname_index" ON "member" USING btree ("userId","orgname") WHERE "member"."userId" IS NOT NULL;--> statement-breakpoint
CREATE UNIQUE INDEX "member_invite_email_orgname_index" ON "member" USING btree ("invite_email","orgname") WHERE "member"."invite_email" IS NOT NULL;--> statement-breakpoint
CREATE UNIQUE INDEX "pipelineStep_name_pipeline_id_index" ON "pipelineStep" USING btree ("name","pipeline_id");--> statement-breakpoint
CREATE UNIQUE INDEX "run_pipeline_run_id_pipeline_step_id_index" ON "run" USING btree ("pipeline_run_id","pipeline_step_id");