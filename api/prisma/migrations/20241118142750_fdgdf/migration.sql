/*
  Warnings:

  - You are about to drop the `PipelineRun` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `ToolRun` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "Content" DROP CONSTRAINT "Content_producedById_fkey";

-- DropForeignKey
ALTER TABLE "PipelineRun" DROP CONSTRAINT "PipelineRun_orgname_fkey";

-- DropForeignKey
ALTER TABLE "PipelineRun" DROP CONSTRAINT "PipelineRun_pipelineId_fkey";

-- DropForeignKey
ALTER TABLE "ToolRun" DROP CONSTRAINT "ToolRun_orgname_fkey";

-- DropForeignKey
ALTER TABLE "ToolRun" DROP CONSTRAINT "ToolRun_pipelineRunId_fkey";

-- DropForeignKey
ALTER TABLE "ToolRun" DROP CONSTRAINT "ToolRun_pipelineStepId_fkey";

-- DropForeignKey
ALTER TABLE "_ContentConsumedBy" DROP CONSTRAINT "_ContentConsumedBy_B_fkey";

-- DropTable
DROP TABLE "PipelineRun";

-- DropTable
DROP TABLE "ToolRun";

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
    "pipelineId" TEXT NOT NULL,
    "pipelineRunId" TEXT,
    "pipelineStepId" TEXT,

    CONSTRAINT "Run_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_producedById_fkey" FOREIGN KEY ("producedById") REFERENCES "Run"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_pipelineId_fkey" FOREIGN KEY ("pipelineId") REFERENCES "Pipeline"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_pipelineRunId_fkey" FOREIGN KEY ("pipelineRunId") REFERENCES "Run"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_pipelineStepId_fkey" FOREIGN KEY ("pipelineStepId") REFERENCES "PipelineStep"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ContentConsumedBy" ADD CONSTRAINT "_ContentConsumedBy_B_fkey" FOREIGN KEY ("B") REFERENCES "Run"("id") ON DELETE CASCADE ON UPDATE CASCADE;
