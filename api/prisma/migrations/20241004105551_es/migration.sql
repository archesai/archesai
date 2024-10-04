/*
  Warnings:

  - The values [API] on the enum `PlanType` will be removed. If these variants are still used in the database, this will fail.

*/
-- AlterEnum
BEGIN;
CREATE TYPE "PlanType_new" AS ENUM ('FREE', 'BASIC', 'STANDARD', 'PREMIUM');
ALTER TABLE "Organization" ALTER COLUMN "plan" DROP DEFAULT;
ALTER TABLE "Organization" ALTER COLUMN "plan" TYPE "PlanType_new" USING ("plan"::text::"PlanType_new");
ALTER TYPE "PlanType" RENAME TO "PlanType_old";
ALTER TYPE "PlanType_new" RENAME TO "PlanType";
DROP TYPE "PlanType_old";
ALTER TABLE "Organization" ALTER COLUMN "plan" SET DEFAULT 'FREE';
COMMIT;
