/*
  Warnings:

  - You are about to drop the column `questionsRemaining` on the `User` table. All the data in the column will be lost.
  - You are about to drop the column `userSetup` on the `User` table. All the data in the column will be lost.
  - You are about to drop the `_DocumentToThread` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "_DocumentToThread" DROP CONSTRAINT "_DocumentToThread_A_fkey";

-- DropForeignKey
ALTER TABLE "_DocumentToThread" DROP CONSTRAINT "_DocumentToThread_B_fkey";

-- AlterTable
ALTER TABLE "User" DROP COLUMN "questionsRemaining",
DROP COLUMN "userSetup";

-- DropTable
DROP TABLE "_DocumentToThread";
