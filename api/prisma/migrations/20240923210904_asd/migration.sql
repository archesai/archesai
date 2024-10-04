/*
  Warnings:

  - You are about to drop the column `answerSimilarity` on the `Citation` table. All the data in the column will be lost.
  - You are about to drop the column `interactionId` on the `Citation` table. All the data in the column will be lost.
  - You are about to drop the column `questionSimilarity` on the `Citation` table. All the data in the column will be lost.
  - You are about to drop the column `textSegmentId` on the `Citation` table. All the data in the column will be lost.
  - You are about to drop the column `animationId` on the `Job` table. All the data in the column will be lost.
  - You are about to drop the column `documentId` on the `Job` table. All the data in the column will be lost.
  - You are about to drop the column `imageId` on the `Job` table. All the data in the column will be lost.
  - You are about to drop the column `agentId` on the `Thread` table. All the data in the column will be lost.
  - You are about to drop the column `defaultOrgname` on the `User` table. All the data in the column will be lost.
  - You are about to drop the column `uid` on the `User` table. All the data in the column will be lost.
  - You are about to drop the `Agent` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Animation` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Document` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Image` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Interaction` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `TextSegment` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Token` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `_AgentToDocument` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `_AgentToToken` table. If the table is not empty, all the data it contains will be lost.
  - Added the required column `contentId` to the `Citation` table without a default value. This is not possible if the table is not empty.
  - Added the required column `messageId` to the `Citation` table without a default value. This is not possible if the table is not empty.
  - Added the required column `similarity` to the `Citation` table without a default value. This is not possible if the table is not empty.
  - Made the column `inviteAccepted` on table `Member` required. This step will fail if there are existing NULL values in that column.
  - Added the required column `chatbotId` to the `Thread` table without a default value. This is not possible if the table is not empty.

*/
-- CreateExtension
CREATE EXTENSION IF NOT EXISTS "vector";

-- CreateEnum
CREATE TYPE "AuthProviderType" AS ENUM ('LOCAL', 'FIREBASE', 'TWITTER');

-- CreateEnum
CREATE TYPE "ContentType" AS ENUM ('DOCUMENT', 'ANIMATION', 'IMAGE');

-- AlterEnum
ALTER TYPE "PlanType" ADD VALUE 'STANDARD';

-- DropForeignKey
ALTER TABLE "Agent" DROP CONSTRAINT "Agent_orgname_fkey";

-- DropForeignKey
ALTER TABLE "Animation" DROP CONSTRAINT "Animation_orgname_fkey";

-- DropForeignKey
ALTER TABLE "Citation" DROP CONSTRAINT "Citation_interactionId_fkey";

-- DropForeignKey
ALTER TABLE "Citation" DROP CONSTRAINT "Citation_textSegmentId_fkey";

-- DropForeignKey
ALTER TABLE "Document" DROP CONSTRAINT "Document_orgname_fkey";

-- DropForeignKey
ALTER TABLE "Image" DROP CONSTRAINT "Image_orgname_fkey";

-- DropForeignKey
ALTER TABLE "Interaction" DROP CONSTRAINT "Interaction_threadId_fkey";

-- DropForeignKey
ALTER TABLE "Job" DROP CONSTRAINT "Job_animationId_fkey";

-- DropForeignKey
ALTER TABLE "Job" DROP CONSTRAINT "Job_documentId_fkey";

-- DropForeignKey
ALTER TABLE "Job" DROP CONSTRAINT "Job_imageId_fkey";

-- DropForeignKey
ALTER TABLE "TextSegment" DROP CONSTRAINT "TextSegment_documentId_fkey";

-- DropForeignKey
ALTER TABLE "TextSegment" DROP CONSTRAINT "TextSegment_orgname_fkey";

-- DropForeignKey
ALTER TABLE "Thread" DROP CONSTRAINT "Thread_agentId_fkey";

-- DropForeignKey
ALTER TABLE "Token" DROP CONSTRAINT "Token_orgname_fkey";

-- DropForeignKey
ALTER TABLE "Token" DROP CONSTRAINT "Token_username_fkey";

-- DropForeignKey
ALTER TABLE "User" DROP CONSTRAINT "User_defaultOrgname_fkey";

-- DropForeignKey
ALTER TABLE "_AgentToDocument" DROP CONSTRAINT "_AgentToDocument_A_fkey";

-- DropForeignKey
ALTER TABLE "_AgentToDocument" DROP CONSTRAINT "_AgentToDocument_B_fkey";

-- DropForeignKey
ALTER TABLE "_AgentToToken" DROP CONSTRAINT "_AgentToToken_A_fkey";

-- DropForeignKey
ALTER TABLE "_AgentToToken" DROP CONSTRAINT "_AgentToToken_B_fkey";

-- DropIndex
DROP INDEX "Job_animationId_key";

-- DropIndex
DROP INDEX "Job_documentId_key";

-- DropIndex
DROP INDEX "Job_imageId_key";

-- DropIndex
DROP INDEX "User_uid_key";

-- AlterTable
ALTER TABLE "Citation" DROP COLUMN "answerSimilarity",
DROP COLUMN "interactionId",
DROP COLUMN "questionSimilarity",
DROP COLUMN "textSegmentId",
ADD COLUMN     "contentId" TEXT NOT NULL,
ADD COLUMN     "messageId" TEXT NOT NULL,
ADD COLUMN     "similarity" DOUBLE PRECISION NOT NULL;

-- AlterTable
ALTER TABLE "Job" DROP COLUMN "animationId",
DROP COLUMN "documentId",
DROP COLUMN "imageId";

-- AlterTable
ALTER TABLE "Member" ALTER COLUMN "inviteAccepted" SET NOT NULL;

-- AlterTable
ALTER TABLE "Thread" DROP COLUMN "agentId",
ADD COLUMN     "chatbotId" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "User" DROP COLUMN "defaultOrgname",
DROP COLUMN "uid",
ADD COLUMN     "defaultOrg" TEXT,
ADD COLUMN     "password" TEXT,
ADD COLUMN     "refreshToken" TEXT;

-- DropTable
DROP TABLE "Agent";

-- DropTable
DROP TABLE "Animation";

-- DropTable
DROP TABLE "Document";

-- DropTable
DROP TABLE "Image";

-- DropTable
DROP TABLE "Interaction";

-- DropTable
DROP TABLE "TextSegment";

-- DropTable
DROP TABLE "Token";

-- DropTable
DROP TABLE "_AgentToDocument";

-- DropTable
DROP TABLE "_AgentToToken";

-- DropEnum
DROP TYPE "AccessScope";

-- DropEnum
DROP TYPE "LLMBase";

-- DropEnum
DROP TYPE "VisibilityStatus";

-- CreateTable
CREATE TABLE "Content" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "text" TEXT,
    "url" TEXT,
    "credits" INTEGER NOT NULL DEFAULT 0,
    "mimeType" TEXT NOT NULL,
    "type" "ContentType" NOT NULL,
    "previewImage" TEXT,
    "orgname" TEXT NOT NULL,
    "jobId" TEXT,
    "buildArgs" JSONB NOT NULL,
    "annotations" JSONB NOT NULL,

    CONSTRAINT "Content_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Chatbot" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL DEFAULT 'Default Search Chatbot',
    "llmBase" TEXT NOT NULL,
    "description" TEXT DEFAULT 'You are an AI-powered search chatbot. You can answer questions about documents.',
    "orgname" TEXT NOT NULL,

    CONSTRAINT "Chatbot_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Message" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "question" TEXT NOT NULL,
    "answer" TEXT NOT NULL,
    "contextLength" INTEGER NOT NULL,
    "answerLength" INTEGER NOT NULL,
    "topK" INTEGER NOT NULL,
    "similarityCutoff" DOUBLE PRECISION NOT NULL DEFAULT 0.7,
    "temperature" DOUBLE PRECISION NOT NULL DEFAULT 0,
    "credits" INTEGER NOT NULL DEFAULT 0,
    "threadId" TEXT NOT NULL,

    CONSTRAINT "Message_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "AuthProvider" (
    "id" TEXT NOT NULL,
    "provider" "AuthProviderType" NOT NULL,
    "providerId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "AuthProvider_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "VectorRecord" (
    "id" TEXT NOT NULL,
    "orgname" TEXT NOT NULL,
    "embedding" vector NOT NULL,
    "contentId" TEXT NOT NULL,
    "text" TEXT NOT NULL,

    CONSTRAINT "VectorRecord_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ApiToken" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL DEFAULT '',
    "role" "RoleType" NOT NULL DEFAULT 'USER',
    "key" TEXT NOT NULL,
    "domains" TEXT NOT NULL DEFAULT '*',
    "orgname" TEXT NOT NULL,
    "username" TEXT,

    CONSTRAINT "ApiToken_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "_ApiTokenToChatbot" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "Content_jobId_key" ON "Content"("jobId");

-- CreateIndex
CREATE UNIQUE INDEX "AuthProvider_provider_providerId_key" ON "AuthProvider"("provider", "providerId");

-- CreateIndex
CREATE UNIQUE INDEX "_ApiTokenToChatbot_AB_unique" ON "_ApiTokenToChatbot"("A", "B");

-- CreateIndex
CREATE INDEX "_ApiTokenToChatbot_B_index" ON "_ApiTokenToChatbot"("B");

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_jobId_fkey" FOREIGN KEY ("jobId") REFERENCES "Job"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Chatbot" ADD CONSTRAINT "Chatbot_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Thread" ADD CONSTRAINT "Thread_chatbotId_fkey" FOREIGN KEY ("chatbotId") REFERENCES "Chatbot"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Message" ADD CONSTRAINT "Message_threadId_fkey" FOREIGN KEY ("threadId") REFERENCES "Thread"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Citation" ADD CONSTRAINT "Citation_messageId_fkey" FOREIGN KEY ("messageId") REFERENCES "Message"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Citation" ADD CONSTRAINT "Citation_contentId_fkey" FOREIGN KEY ("contentId") REFERENCES "Content"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "AuthProvider" ADD CONSTRAINT "AuthProvider_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "VectorRecord" ADD CONSTRAINT "VectorRecord_contentId_fkey" FOREIGN KEY ("contentId") REFERENCES "Content"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ApiToken" ADD CONSTRAINT "ApiToken_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ApiToken" ADD CONSTRAINT "ApiToken_username_fkey" FOREIGN KEY ("username") REFERENCES "User"("username") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ApiTokenToChatbot" ADD CONSTRAINT "_ApiTokenToChatbot_A_fkey" FOREIGN KEY ("A") REFERENCES "ApiToken"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ApiTokenToChatbot" ADD CONSTRAINT "_ApiTokenToChatbot_B_fkey" FOREIGN KEY ("B") REFERENCES "Chatbot"("id") ON DELETE CASCADE ON UPDATE CASCADE;
