/*
  Warnings:

  - You are about to drop the column `threadId` on the `PipelineRun` table. All the data in the column will be lost.
  - You are about to drop the `Thread` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "PipelineRun" DROP CONSTRAINT "PipelineRun_threadId_fkey";

-- DropForeignKey
ALTER TABLE "Thread" DROP CONSTRAINT "Thread_orgname_fkey";

-- AlterTable
ALTER TABLE "PipelineRun" DROP COLUMN "threadId";

-- DropTable
DROP TABLE "Thread";

-- CreateTable
CREATE TABLE "Labels" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "orgname" TEXT NOT NULL,

    CONSTRAINT "Labels_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "Labels_name_orgname_key" ON "Labels"("name", "orgname");

-- AddForeignKey
ALTER TABLE "Labels" ADD CONSTRAINT "Labels_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;
