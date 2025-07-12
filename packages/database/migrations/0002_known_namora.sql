ALTER TABLE "organizations" DROP CONSTRAINT "organizations_orgname_unique";--> statement-breakpoint
ALTER TABLE "artifacts" ADD COLUMN "name" text;--> statement-breakpoint
ALTER TABLE "organizations" ADD CONSTRAINT "organizations_organizationId_unique" UNIQUE("organizationId");