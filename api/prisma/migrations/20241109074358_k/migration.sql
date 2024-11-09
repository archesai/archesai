-- CreateExtension
CREATE EXTENSION IF NOT EXISTS "vector";

-- CreateEnum
CREATE TYPE "RoleType" AS ENUM ('USER', 'ADMIN');

-- CreateEnum
CREATE TYPE "AuthProviderType" AS ENUM ('LOCAL', 'FIREBASE', 'TWITTER');

-- CreateEnum
CREATE TYPE "ToolIOType" AS ENUM ('TEXT', 'IMAGE', 'VIDEO', 'AUDIO');

-- CreateEnum
CREATE TYPE "RunRole" AS ENUM ('INPUT', 'OUTPUT');

-- CreateEnum
CREATE TYPE "RunStatus" AS ENUM ('QUEUED', 'PROCESSING', 'COMPLETE', 'ERROR');

-- CreateEnum
CREATE TYPE "PlanType" AS ENUM ('FREE', 'BASIC', 'STANDARD', 'PREMIUM', 'UNLIMITED');

-- CreateEnum
CREATE TYPE "ARTokenType" AS ENUM ('EMAIL_VERIFICATION', 'PASSWORD_RESET', 'EMAIL_CHANGE');

-- CreateTable
CREATE TABLE "Content" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "url" TEXT,
    "text" TEXT,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "embedding" vector,
    "credits" INTEGER NOT NULL DEFAULT 0,
    "mimeType" TEXT,
    "previewImage" TEXT,
    "orgname" TEXT NOT NULL,
    "parentId" TEXT,

    CONSTRAINT "Content_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "RunContent" (
    "id" TEXT NOT NULL,
    "role" "RunRole" NOT NULL,
    "pipelineRunId" TEXT NOT NULL,
    "contentId" TEXT NOT NULL,
    "transformationId" TEXT,

    CONSTRAINT "RunContent_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "PipelineRun" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "status" "RunStatus" NOT NULL DEFAULT 'QUEUED',
    "startedAt" TIMESTAMP(3),
    "completedAt" TIMESTAMP(3),
    "progress" DOUBLE PRECISION NOT NULL DEFAULT 0,
    "error" TEXT,
    "orgname" TEXT NOT NULL,
    "pipelineId" TEXT NOT NULL,
    "threadId" TEXT,

    CONSTRAINT "PipelineRun_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Transformation" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "status" "RunStatus" NOT NULL DEFAULT 'QUEUED',
    "startedAt" TIMESTAMP(3),
    "completedAt" TIMESTAMP(3),
    "progress" DOUBLE PRECISION NOT NULL DEFAULT 0,
    "error" TEXT,
    "pipelineRunId" TEXT NOT NULL,
    "pipelineStepId" TEXT NOT NULL,

    CONSTRAINT "Transformation_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Tool" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "inputType" "ToolIOType" NOT NULL,
    "outputType" "ToolIOType" NOT NULL,
    "toolBase" TEXT NOT NULL,
    "orgname" TEXT NOT NULL,

    CONSTRAINT "Tool_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Pipeline" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "orgname" TEXT NOT NULL,

    CONSTRAINT "Pipeline_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "PipelineStep" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "pipelineId" TEXT NOT NULL,
    "toolId" TEXT NOT NULL,
    "dependsOnId" TEXT,

    CONSTRAINT "PipelineStep_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Thread" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL DEFAULT 'New Thread',
    "credits" INTEGER NOT NULL DEFAULT 0,
    "orgname" TEXT NOT NULL,

    CONSTRAINT "Thread_pkey" PRIMARY KEY ("id")
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
CREATE TABLE "User" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "username" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "firstName" TEXT NOT NULL DEFAULT '',
    "lastName" TEXT NOT NULL DEFAULT '',
    "emailVerified" BOOLEAN NOT NULL DEFAULT false,
    "deactivated" BOOLEAN NOT NULL DEFAULT false,
    "photoUrl" TEXT,
    "defaultOrgname" TEXT,
    "password" TEXT,
    "refreshToken" TEXT,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Organization" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "stripeCustomerId" TEXT NOT NULL,
    "orgname" TEXT NOT NULL,
    "billingEmail" TEXT NOT NULL,
    "credits" INTEGER NOT NULL DEFAULT 0,
    "plan" "PlanType" NOT NULL DEFAULT 'FREE',

    CONSTRAINT "Organization_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Member" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "inviteAccepted" BOOLEAN NOT NULL DEFAULT false,
    "inviteEmail" TEXT NOT NULL,
    "role" "RoleType" NOT NULL DEFAULT 'USER',
    "orgname" TEXT NOT NULL,
    "username" TEXT,

    CONSTRAINT "Member_pkey" PRIMARY KEY ("id")
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
CREATE TABLE "ARToken" (
    "id" TEXT NOT NULL,
    "type" "ARTokenType" NOT NULL,
    "token" TEXT NOT NULL,
    "expiresAt" TIMESTAMP(3) NOT NULL,
    "userId" TEXT NOT NULL,
    "newEmail" TEXT,

    CONSTRAINT "ARToken_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "RunContent_pipelineRunId_idx" ON "RunContent"("pipelineRunId");

-- CreateIndex
CREATE INDEX "RunContent_contentId_idx" ON "RunContent"("contentId");

-- CreateIndex
CREATE UNIQUE INDEX "RunContent_pipelineRunId_contentId_role_key" ON "RunContent"("pipelineRunId", "contentId", "role");

-- CreateIndex
CREATE UNIQUE INDEX "AuthProvider_provider_providerId_key" ON "AuthProvider"("provider", "providerId");

-- CreateIndex
CREATE UNIQUE INDEX "User_username_key" ON "User"("username");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "Organization_stripeCustomerId_key" ON "Organization"("stripeCustomerId");

-- CreateIndex
CREATE UNIQUE INDEX "Organization_orgname_key" ON "Organization"("orgname");

-- CreateIndex
CREATE UNIQUE INDEX "Member_inviteEmail_orgname_key" ON "Member"("inviteEmail", "orgname");

-- CreateIndex
CREATE UNIQUE INDEX "Member_username_orgname_key" ON "Member"("username", "orgname");

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_parentId_fkey" FOREIGN KEY ("parentId") REFERENCES "Content"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "RunContent" ADD CONSTRAINT "RunContent_pipelineRunId_fkey" FOREIGN KEY ("pipelineRunId") REFERENCES "PipelineRun"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "RunContent" ADD CONSTRAINT "RunContent_contentId_fkey" FOREIGN KEY ("contentId") REFERENCES "Content"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "RunContent" ADD CONSTRAINT "RunContent_transformationId_fkey" FOREIGN KEY ("transformationId") REFERENCES "Transformation"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "PipelineRun" ADD CONSTRAINT "PipelineRun_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "PipelineRun" ADD CONSTRAINT "PipelineRun_pipelineId_fkey" FOREIGN KEY ("pipelineId") REFERENCES "Pipeline"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "PipelineRun" ADD CONSTRAINT "PipelineRun_threadId_fkey" FOREIGN KEY ("threadId") REFERENCES "Thread"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Transformation" ADD CONSTRAINT "Transformation_pipelineRunId_fkey" FOREIGN KEY ("pipelineRunId") REFERENCES "PipelineRun"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Transformation" ADD CONSTRAINT "Transformation_pipelineStepId_fkey" FOREIGN KEY ("pipelineStepId") REFERENCES "PipelineStep"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Tool" ADD CONSTRAINT "Tool_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Pipeline" ADD CONSTRAINT "Pipeline_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "PipelineStep" ADD CONSTRAINT "PipelineStep_pipelineId_fkey" FOREIGN KEY ("pipelineId") REFERENCES "Pipeline"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "PipelineStep" ADD CONSTRAINT "PipelineStep_toolId_fkey" FOREIGN KEY ("toolId") REFERENCES "Tool"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "PipelineStep" ADD CONSTRAINT "PipelineStep_dependsOnId_fkey" FOREIGN KEY ("dependsOnId") REFERENCES "PipelineStep"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Thread" ADD CONSTRAINT "Thread_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "AuthProvider" ADD CONSTRAINT "AuthProvider_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Member" ADD CONSTRAINT "Member_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Member" ADD CONSTRAINT "Member_username_fkey" FOREIGN KEY ("username") REFERENCES "User"("username") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ApiToken" ADD CONSTRAINT "ApiToken_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ApiToken" ADD CONSTRAINT "ApiToken_username_fkey" FOREIGN KEY ("username") REFERENCES "User"("username") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ARToken" ADD CONSTRAINT "ARToken_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;
