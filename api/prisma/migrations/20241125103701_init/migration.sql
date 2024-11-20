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
CREATE TYPE "RunType" AS ENUM ('PIPELINE_RUN', 'TOOL_RUN');

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
    "producedById" TEXT,

    CONSTRAINT "Content_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Run" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL DEFAULT 'New Run',
    "status" "RunStatus" NOT NULL DEFAULT 'QUEUED',
    "startedAt" TIMESTAMP(3),
    "completedAt" TIMESTAMP(3),
    "progress" DOUBLE PRECISION NOT NULL DEFAULT 0,
    "error" TEXT,
    "orgname" TEXT NOT NULL,
    "runType" "RunType" NOT NULL,
    "pipelineId" TEXT,
    "toolId" TEXT,
    "pipelineRunId" TEXT,
    "pipelineStepId" TEXT,

    CONSTRAINT "Run_pkey" PRIMARY KEY ("id")
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
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "pipelineId" TEXT NOT NULL,
    "toolId" TEXT NOT NULL,

    CONSTRAINT "PipelineStep_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Label" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "orgname" TEXT NOT NULL,

    CONSTRAINT "Label_pkey" PRIMARY KEY ("id")
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

-- CreateTable
CREATE TABLE "_ContentConsumedBy" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateTable
CREATE TABLE "_ContentToLabel" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateTable
CREATE TABLE "_PipelineStepDependencies" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "Run_pipelineRunId_pipelineStepId_key" ON "Run"("pipelineRunId", "pipelineStepId");

-- CreateIndex
CREATE UNIQUE INDEX "PipelineStep_name_pipelineId_key" ON "PipelineStep"("name", "pipelineId");

-- CreateIndex
CREATE UNIQUE INDEX "Label_name_orgname_key" ON "Label"("name", "orgname");

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

-- CreateIndex
CREATE UNIQUE INDEX "_ContentConsumedBy_AB_unique" ON "_ContentConsumedBy"("A", "B");

-- CreateIndex
CREATE INDEX "_ContentConsumedBy_B_index" ON "_ContentConsumedBy"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_ContentToLabel_AB_unique" ON "_ContentToLabel"("A", "B");

-- CreateIndex
CREATE INDEX "_ContentToLabel_B_index" ON "_ContentToLabel"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_PipelineStepDependencies_AB_unique" ON "_PipelineStepDependencies"("A", "B");

-- CreateIndex
CREATE INDEX "_PipelineStepDependencies_B_index" ON "_PipelineStepDependencies"("B");

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_parentId_fkey" FOREIGN KEY ("parentId") REFERENCES "Content"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_producedById_fkey" FOREIGN KEY ("producedById") REFERENCES "Run"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_pipelineId_fkey" FOREIGN KEY ("pipelineId") REFERENCES "Pipeline"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_toolId_fkey" FOREIGN KEY ("toolId") REFERENCES "Tool"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_pipelineRunId_fkey" FOREIGN KEY ("pipelineRunId") REFERENCES "Run"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_pipelineStepId_fkey" FOREIGN KEY ("pipelineStepId") REFERENCES "PipelineStep"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Tool" ADD CONSTRAINT "Tool_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Pipeline" ADD CONSTRAINT "Pipeline_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "PipelineStep" ADD CONSTRAINT "PipelineStep_pipelineId_fkey" FOREIGN KEY ("pipelineId") REFERENCES "Pipeline"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "PipelineStep" ADD CONSTRAINT "PipelineStep_toolId_fkey" FOREIGN KEY ("toolId") REFERENCES "Tool"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Label" ADD CONSTRAINT "Label_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

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

-- AddForeignKey
ALTER TABLE "_ContentConsumedBy" ADD CONSTRAINT "_ContentConsumedBy_A_fkey" FOREIGN KEY ("A") REFERENCES "Content"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ContentConsumedBy" ADD CONSTRAINT "_ContentConsumedBy_B_fkey" FOREIGN KEY ("B") REFERENCES "Run"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ContentToLabel" ADD CONSTRAINT "_ContentToLabel_A_fkey" FOREIGN KEY ("A") REFERENCES "Content"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ContentToLabel" ADD CONSTRAINT "_ContentToLabel_B_fkey" FOREIGN KEY ("B") REFERENCES "Label"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_PipelineStepDependencies" ADD CONSTRAINT "_PipelineStepDependencies_A_fkey" FOREIGN KEY ("A") REFERENCES "PipelineStep"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_PipelineStepDependencies" ADD CONSTRAINT "_PipelineStepDependencies_B_fkey" FOREIGN KEY ("B") REFERENCES "PipelineStep"("id") ON DELETE CASCADE ON UPDATE CASCADE;
