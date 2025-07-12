CREATE TYPE "public"."authType" AS ENUM('email', 'oauth', 'oidc', 'webauthn');--> statement-breakpoint
CREATE TYPE "public"."planType" AS ENUM('BASIC', 'FREE', 'PREMIUM', 'STANDARD', 'UNLIMITED');--> statement-breakpoint
CREATE TYPE "public"."providerType" AS ENUM('API_KEY', 'FIREBASE', 'LOCAL', 'TWITTER');--> statement-breakpoint
CREATE TYPE "public"."role" AS ENUM('ADMIN', 'USER');--> statement-breakpoint
CREATE TYPE "public"."runStatus" AS ENUM('COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED');--> statement-breakpoint
CREATE TYPE "public"."RunType" AS ENUM('PIPELINE_RUN', 'TOOL_RUN');--> statement-breakpoint
CREATE TYPE "public"."toolIO" AS ENUM('AUDIO', 'IMAGE', 'TEXT', 'VIDEO');--> statement-breakpoint
CREATE TYPE "public"."verificationTokenType" AS ENUM('EMAIL_CHANGE', 'EMAIL_VERIFICATION', 'PASSWORD_RESET');--> statement-breakpoint
CREATE TABLE "accounts" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"accessToken" text,
	"accessTokenExpiresAt" timestamp,
	"accountId" text NOT NULL,
	"idToken" text,
	"password" text,
	"providerId" text NOT NULL,
	"refreshToken" text,
	"refreshTokenExpiresAt" timestamp,
	"scope" text,
	"userId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "api-tokens" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"key" text NOT NULL,
	"organizationId" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL
);
--> statement-breakpoint
CREATE TABLE "artifacts" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"credits" integer DEFAULT 0 NOT NULL,
	"description" text,
	"mimeType" text,
	"organizationId" text NOT NULL,
	"parentId" text,
	"previewImage" text,
	"producerId" text,
	"text" text,
	"url" text
);
--> statement-breakpoint
CREATE TABLE "invitations" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"accepted" boolean DEFAULT false,
	"email" text NOT NULL,
	"organizationId" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL
);
--> statement-breakpoint
CREATE TABLE "labels" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"name" text NOT NULL,
	"organizationId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "labelToArtifact" (
	"artifactId" text NOT NULL,
	"labelId" text NOT NULL,
	CONSTRAINT "labelToArtifact_labelId_artifactId_pk" PRIMARY KEY("labelId","artifactId")
);
--> statement-breakpoint
CREATE TABLE "members" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"invitationId" text,
	"organizationId" text NOT NULL,
	"role" "role" DEFAULT 'USER' NOT NULL,
	"userId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "organizations" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"billingEmail" text NOT NULL,
	"credits" integer DEFAULT 0 NOT NULL,
	"organizationId" text NOT NULL,
	"plan" "planType" DEFAULT 'FREE' NOT NULL,
	"stripeCustomerId" text,
	CONSTRAINT "organizations_orgname_unique" UNIQUE("organizationId"),
	CONSTRAINT "organizations_stripeCustomerId_unique" UNIQUE("stripeCustomerId")
);
--> statement-breakpoint
CREATE TABLE "pipelines" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"description" text,
	"organizationId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "pipeline-steps" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"pipelineId" text NOT NULL,
	"toolId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "pipelineStepToDependency" (
	"pipelineStepId" text NOT NULL,
	"prerequisiteId" text NOT NULL,
	CONSTRAINT "pipelineStepToDependency_pipelineStepId_prerequisiteId_pk" PRIMARY KEY("pipelineStepId","prerequisiteId")
);
--> statement-breakpoint
CREATE TABLE "runs" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"completedAt" timestamp(3),
	"error" text,
	"organizationId" text NOT NULL,
	"pipelineId" text,
	"progress" double precision DEFAULT 0 NOT NULL,
	"runType" "RunType" NOT NULL,
	"startedAt" timestamp(3),
	"status" "runStatus" DEFAULT 'QUEUED' NOT NULL,
	"toolId" text
);
--> statement-breakpoint
CREATE TABLE "session" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"expiresAt" timestamp NOT NULL,
	"ipAddress" text,
	"token" text NOT NULL,
	"userAgent" text,
	"userId" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "tools" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"description" text NOT NULL,
	"inputType" "toolIO" NOT NULL,
	"name" text NOT NULL,
	"organizationId" text NOT NULL,
	"outputType" "toolIO" NOT NULL,
	"toolBase" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "users" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"deactivated" boolean DEFAULT false NOT NULL,
	"email" text NOT NULL,
	"emailVerified" boolean DEFAULT false NOT NULL,
	"image" text,
	"name" text NOT NULL,
	CONSTRAINT "users_email_unique" UNIQUE("email")
);
--> statement-breakpoint
CREATE TABLE "verification-tokens" (
	"createdAt" timestamp DEFAULT now() NOT NULL,
	"id" text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"updatedAt" timestamp DEFAULT now() NOT NULL,
	"expiresAt" timestamp NOT NULL,
	"identifier" text NOT NULL,
	"value" text NOT NULL
);
--> statement-breakpoint
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "api-tokens" ADD CONSTRAINT "api-tokens_organizationId_organizations_id_fk" FOREIGN KEY ("organizationId") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_organizationId_organizations_id_fk" FOREIGN KEY ("organizationId") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_parentId_artifacts_id_fk" FOREIGN KEY ("parentId") REFERENCES "public"."artifacts"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "artifacts" ADD CONSTRAINT "artifacts_producerId_runs_id_fk" FOREIGN KEY ("producerId") REFERENCES "public"."runs"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "invitations" ADD CONSTRAINT "invitations_organizationId_organizations_id_fk" FOREIGN KEY ("organizationId") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "labels" ADD CONSTRAINT "labels_organizationId_organizations_id_fk" FOREIGN KEY ("organizationId") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "labelToArtifact" ADD CONSTRAINT "labelToArtifact_artifactId_artifacts_id_fk" FOREIGN KEY ("artifactId") REFERENCES "public"."artifacts"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "labelToArtifact" ADD CONSTRAINT "labelToArtifact_labelId_labels_id_fk" FOREIGN KEY ("labelId") REFERENCES "public"."labels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_invitationId_invitations_id_fk" FOREIGN KEY ("invitationId") REFERENCES "public"."invitations"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_organizationId_organizations_id_fk" FOREIGN KEY ("organizationId") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "members" ADD CONSTRAINT "members_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "pipelines" ADD CONSTRAINT "pipelines_organizationId_organizations_id_fk" FOREIGN KEY ("organizationId") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipeline-steps" ADD CONSTRAINT "pipeline-steps_pipelineId_pipelines_id_fk" FOREIGN KEY ("pipelineId") REFERENCES "public"."pipelines"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipeline-steps" ADD CONSTRAINT "pipeline-steps_toolId_tools_id_fk" FOREIGN KEY ("toolId") REFERENCES "public"."tools"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "pipelineStepToDependency" ADD CONSTRAINT "pipelineStepToDependency_pipelineStepId_pipeline-steps_id_fk" FOREIGN KEY ("pipelineStepId") REFERENCES "public"."pipeline-steps"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "pipelineStepToDependency" ADD CONSTRAINT "pipelineStepToDependency_prerequisiteId_pipeline-steps_id_fk" FOREIGN KEY ("prerequisiteId") REFERENCES "public"."pipeline-steps"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_organizationId_organizations_id_fk" FOREIGN KEY ("organizationId") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_pipelineId_pipelines_id_fk" FOREIGN KEY ("pipelineId") REFERENCES "public"."pipelines"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "runs" ADD CONSTRAINT "runs_toolId_tools_id_fk" FOREIGN KEY ("toolId") REFERENCES "public"."tools"("id") ON DELETE set null ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "session" ADD CONSTRAINT "session_userId_users_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "tools" ADD CONSTRAINT "tools_organizationId_organizations_id_fk" FOREIGN KEY ("organizationId") REFERENCES "public"."organizations"("id") ON DELETE cascade ON UPDATE cascade;--> statement-breakpoint
CREATE UNIQUE INDEX "labels_name_organizationId_index" ON "labels" USING btree ("name","organizationId");