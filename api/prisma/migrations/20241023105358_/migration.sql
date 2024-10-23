/*
  Warnings:

  - You are about to drop the column `jobType` on the `Job` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "Job" DROP COLUMN "jobType",
ADD COLUMN     "error" TEXT;
