/*
  Warnings:

  - You are about to drop the column `defaultOrg` on the `User` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "User" DROP COLUMN "defaultOrg",
ADD COLUMN     "defaultOrgname" TEXT;
