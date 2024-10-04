-- CreateEnum
CREATE TYPE "RoleType" AS ENUM ('USER', 'ADMIN');

-- CreateEnum
CREATE TYPE "JobStatus" AS ENUM ('QUEUED', 'PROCESSING', 'COMPLETE', 'ERROR');

-- CreateEnum
CREATE TYPE "VisibilityStatus" AS ENUM ('PUBLIC', 'PRIVATE');

-- CreateEnum
CREATE TYPE "PlanType" AS ENUM ('FREE', 'BASIC', 'PREMIUM', 'API');

-- CreateEnum
CREATE TYPE "LLMBase" AS ENUM ('GPT_3_5_TURBO_16_K', 'GPT_4', 'MISTRAL_7B');

-- CreateEnum
CREATE TYPE "AccessScope" AS ENUM ('ORGANIZATION', 'DOCUMENT');

-- CreateEnum
CREATE TYPE "ToolDataType" AS ENUM ('TEXT', 'IMAGE', 'AUDIO', 'VIDEO');

-- CreateTable
CREATE TABLE "User" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "username" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "uid" TEXT NOT NULL,
    "firstName" TEXT NOT NULL DEFAULT '',
    "lastName" TEXT NOT NULL DEFAULT '',
    "userSetup" BOOLEAN NOT NULL DEFAULT false,
    "emailVerified" BOOLEAN NOT NULL DEFAULT false,
    "defaultOrgname" TEXT,
    "deactivated" BOOLEAN NOT NULL DEFAULT false,
    "questionsRemaining" INTEGER NOT NULL,
    "photoUrl" TEXT,

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
    "inviteAccepted" BOOLEAN DEFAULT false,
    "inviteEmail" TEXT NOT NULL,
    "role" "RoleType" NOT NULL DEFAULT 'USER',
    "orgname" TEXT NOT NULL,
    "username" TEXT,

    CONSTRAINT "Member_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Document" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "documentUrl" TEXT NOT NULL,
    "visibility" "VisibilityStatus" NOT NULL,
    "summary" TEXT,
    "orgname" TEXT NOT NULL,
    "credits" INTEGER NOT NULL DEFAULT 0,
    "chunkSize" INTEGER NOT NULL DEFAULT 200,
    "delimiter" TEXT,
    "contentType" TEXT NOT NULL,
    "previewImage" TEXT,

    CONSTRAINT "Document_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Animation" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "animationUrl" TEXT NOT NULL,
    "visibility" "VisibilityStatus" NOT NULL,
    "orgname" TEXT NOT NULL,
    "credits" INTEGER NOT NULL DEFAULT 0,
    "contentType" TEXT NOT NULL,
    "width" INTEGER NOT NULL,
    "height" INTEGER NOT NULL,
    "useInit" BOOLEAN NOT NULL,
    "initImage" TEXT NOT NULL,
    "useAudio" BOOLEAN NOT NULL,
    "audioUrl" TEXT NOT NULL,
    "audioStart" DOUBLE PRECISION NOT NULL,
    "animationPrompts" TEXT NOT NULL,
    "maxFrames" INTEGER NOT NULL,
    "fps" INTEGER NOT NULL,
    "translationX" TEXT NOT NULL,
    "translationY" TEXT NOT NULL,
    "translationZ" TEXT NOT NULL,
    "translation3DX" TEXT NOT NULL,
    "translation3DY" TEXT NOT NULL,
    "translation3DZ" TEXT NOT NULL,
    "zoom" TEXT NOT NULL,
    "mode" TEXT NOT NULL,
    "border" TEXT NOT NULL,
    "previewImage" TEXT,

    CONSTRAINT "Animation_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Image" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "imageUrl" TEXT NOT NULL,
    "visibility" "VisibilityStatus" NOT NULL,
    "orgname" TEXT NOT NULL,
    "credits" INTEGER NOT NULL DEFAULT 0,
    "contentType" TEXT NOT NULL,
    "width" INTEGER NOT NULL,
    "height" INTEGER NOT NULL,
    "useInit" BOOLEAN NOT NULL,
    "initImage" TEXT NOT NULL,
    "prompt" TEXT NOT NULL,

    CONSTRAINT "Image_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "TextSegment" (
    "id" TEXT NOT NULL,
    "page" INTEGER NOT NULL,
    "index" INTEGER NOT NULL,
    "text" TEXT NOT NULL,
    "documentId" TEXT,
    "orgname" TEXT NOT NULL,
    "vectorDbId" TEXT NOT NULL,

    CONSTRAINT "TextSegment_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Token" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL DEFAULT '',
    "role" "RoleType" NOT NULL DEFAULT 'USER',
    "snippet" TEXT NOT NULL,
    "orgname" TEXT NOT NULL,
    "username" TEXT,
    "domains" TEXT NOT NULL DEFAULT '*',

    CONSTRAINT "Token_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Thread" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "orgname" TEXT NOT NULL,
    "agentId" TEXT NOT NULL,
    "name" TEXT NOT NULL DEFAULT 'New Thread',
    "credits" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "Thread_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Agent" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "orgname" TEXT NOT NULL,
    "accessScope" "AccessScope" NOT NULL DEFAULT 'ORGANIZATION',
    "name" TEXT NOT NULL DEFAULT 'Default Search Agent',
    "llmBase" "LLMBase" NOT NULL DEFAULT 'GPT_3_5_TURBO_16_K',
    "description" TEXT DEFAULT 'You are an AI-powered search agent. You can answer questions about documents.',

    CONSTRAINT "Agent_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Interaction" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "threadId" TEXT NOT NULL,
    "question" TEXT NOT NULL,
    "answer" TEXT NOT NULL,
    "contextLength" INTEGER NOT NULL,
    "answerLength" INTEGER NOT NULL,
    "topK" INTEGER NOT NULL,
    "similarityCutoff" DOUBLE PRECISION NOT NULL DEFAULT 0.7,
    "temperature" DOUBLE PRECISION NOT NULL DEFAULT 0,
    "credits" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "Interaction_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ToolUsage" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "toolId" TEXT NOT NULL,
    "interactionId" TEXT NOT NULL,

    CONSTRAINT "ToolUsage_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Citation" (
    "id" TEXT NOT NULL,
    "answerSimilarity" DOUBLE PRECISION NOT NULL,
    "questionSimilarity" DOUBLE PRECISION NOT NULL,
    "interactionId" TEXT NOT NULL,
    "textSegmentId" TEXT NOT NULL,

    CONSTRAINT "Citation_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Tool" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "inputs" "ToolDataType"[],
    "outputs" "ToolDataType"[],
    "orgname" TEXT,

    CONSTRAINT "Tool_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Job" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "startedAt" TIMESTAMP(3),
    "completedAt" TIMESTAMP(3),
    "status" "JobStatus" NOT NULL DEFAULT 'QUEUED',
    "progress" DOUBLE PRECISION DEFAULT 0,
    "jobType" TEXT NOT NULL,
    "orgname" TEXT NOT NULL,
    "documentId" TEXT,
    "animationId" TEXT,
    "imageId" TEXT,

    CONSTRAINT "Job_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "_DocumentToThread" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateTable
CREATE TABLE "_AgentToDocument" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateTable
CREATE TABLE "_AgentToToken" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateTable
CREATE TABLE "_AgentToTool" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "User_username_key" ON "User"("username");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "User_uid_key" ON "User"("uid");

-- CreateIndex
CREATE UNIQUE INDEX "Organization_stripeCustomerId_key" ON "Organization"("stripeCustomerId");

-- CreateIndex
CREATE UNIQUE INDEX "Organization_orgname_key" ON "Organization"("orgname");

-- CreateIndex
CREATE UNIQUE INDEX "Member_inviteEmail_orgname_key" ON "Member"("inviteEmail", "orgname");

-- CreateIndex
CREATE UNIQUE INDEX "Member_username_orgname_key" ON "Member"("username", "orgname");

-- CreateIndex
CREATE UNIQUE INDEX "TextSegment_vectorDbId_documentId_key" ON "TextSegment"("vectorDbId", "documentId");

-- CreateIndex
CREATE UNIQUE INDEX "Job_documentId_key" ON "Job"("documentId");

-- CreateIndex
CREATE UNIQUE INDEX "Job_animationId_key" ON "Job"("animationId");

-- CreateIndex
CREATE UNIQUE INDEX "Job_imageId_key" ON "Job"("imageId");

-- CreateIndex
CREATE UNIQUE INDEX "_DocumentToThread_AB_unique" ON "_DocumentToThread"("A", "B");

-- CreateIndex
CREATE INDEX "_DocumentToThread_B_index" ON "_DocumentToThread"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_AgentToDocument_AB_unique" ON "_AgentToDocument"("A", "B");

-- CreateIndex
CREATE INDEX "_AgentToDocument_B_index" ON "_AgentToDocument"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_AgentToToken_AB_unique" ON "_AgentToToken"("A", "B");

-- CreateIndex
CREATE INDEX "_AgentToToken_B_index" ON "_AgentToToken"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_AgentToTool_AB_unique" ON "_AgentToTool"("A", "B");

-- CreateIndex
CREATE INDEX "_AgentToTool_B_index" ON "_AgentToTool"("B");

-- AddForeignKey
ALTER TABLE "User" ADD CONSTRAINT "User_defaultOrgname_fkey" FOREIGN KEY ("defaultOrgname") REFERENCES "Organization"("orgname") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Member" ADD CONSTRAINT "Member_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Member" ADD CONSTRAINT "Member_username_fkey" FOREIGN KEY ("username") REFERENCES "User"("username") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Document" ADD CONSTRAINT "Document_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Animation" ADD CONSTRAINT "Animation_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Image" ADD CONSTRAINT "Image_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "TextSegment" ADD CONSTRAINT "TextSegment_documentId_fkey" FOREIGN KEY ("documentId") REFERENCES "Document"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "TextSegment" ADD CONSTRAINT "TextSegment_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Token" ADD CONSTRAINT "Token_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Token" ADD CONSTRAINT "Token_username_fkey" FOREIGN KEY ("username") REFERENCES "User"("username") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Thread" ADD CONSTRAINT "Thread_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Thread" ADD CONSTRAINT "Thread_agentId_fkey" FOREIGN KEY ("agentId") REFERENCES "Agent"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Agent" ADD CONSTRAINT "Agent_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Interaction" ADD CONSTRAINT "Interaction_threadId_fkey" FOREIGN KEY ("threadId") REFERENCES "Thread"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ToolUsage" ADD CONSTRAINT "ToolUsage_toolId_fkey" FOREIGN KEY ("toolId") REFERENCES "Tool"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ToolUsage" ADD CONSTRAINT "ToolUsage_interactionId_fkey" FOREIGN KEY ("interactionId") REFERENCES "Interaction"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Citation" ADD CONSTRAINT "Citation_interactionId_fkey" FOREIGN KEY ("interactionId") REFERENCES "Interaction"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Citation" ADD CONSTRAINT "Citation_textSegmentId_fkey" FOREIGN KEY ("textSegmentId") REFERENCES "TextSegment"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Tool" ADD CONSTRAINT "Tool_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Job" ADD CONSTRAINT "Job_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Job" ADD CONSTRAINT "Job_documentId_fkey" FOREIGN KEY ("documentId") REFERENCES "Document"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Job" ADD CONSTRAINT "Job_animationId_fkey" FOREIGN KEY ("animationId") REFERENCES "Animation"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Job" ADD CONSTRAINT "Job_imageId_fkey" FOREIGN KEY ("imageId") REFERENCES "Image"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_DocumentToThread" ADD CONSTRAINT "_DocumentToThread_A_fkey" FOREIGN KEY ("A") REFERENCES "Document"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_DocumentToThread" ADD CONSTRAINT "_DocumentToThread_B_fkey" FOREIGN KEY ("B") REFERENCES "Thread"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_AgentToDocument" ADD CONSTRAINT "_AgentToDocument_A_fkey" FOREIGN KEY ("A") REFERENCES "Agent"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_AgentToDocument" ADD CONSTRAINT "_AgentToDocument_B_fkey" FOREIGN KEY ("B") REFERENCES "Document"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_AgentToToken" ADD CONSTRAINT "_AgentToToken_A_fkey" FOREIGN KEY ("A") REFERENCES "Agent"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_AgentToToken" ADD CONSTRAINT "_AgentToToken_B_fkey" FOREIGN KEY ("B") REFERENCES "Token"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_AgentToTool" ADD CONSTRAINT "_AgentToTool_A_fkey" FOREIGN KEY ("A") REFERENCES "Agent"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_AgentToTool" ADD CONSTRAINT "_AgentToTool_B_fkey" FOREIGN KEY ("B") REFERENCES "Tool"("id") ON DELETE CASCADE ON UPDATE CASCADE;
