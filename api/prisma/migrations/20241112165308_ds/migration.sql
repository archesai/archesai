/*
  Warnings:

  - You are about to drop the column `dependsOnId` on the `PipelineStep` table. All the data in the column will be lost.
  - You are about to drop the `Transformation` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "Content" DROP CONSTRAINT "Content_producedById_fkey";

-- DropForeignKey
ALTER TABLE "PipelineStep" DROP CONSTRAINT "PipelineStep_dependsOnId_fkey";

-- DropForeignKey
ALTER TABLE "Transformation" DROP CONSTRAINT "Transformation_pipelineRunId_fkey";

-- DropForeignKey
ALTER TABLE "Transformation" DROP CONSTRAINT "Transformation_pipelineStepId_fkey";

-- DropForeignKey
ALTER TABLE "_ContentConsumedBy" DROP CONSTRAINT "_ContentConsumedBy_B_fkey";

-- AlterTable
ALTER TABLE "PipelineStep" DROP COLUMN "dependsOnId";

-- DropTable
DROP TABLE "Transformation";

-- CreateTable
CREATE TABLE "ToolRun" (
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
    "pipelineRunId" TEXT NOT NULL,
    "pipelineStepId" TEXT NOT NULL,

    CONSTRAINT "ToolRun_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "_PipelineStepDependencies" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "ToolRun_pipelineRunId_pipelineStepId_key" ON "ToolRun"("pipelineRunId", "pipelineStepId");

-- CreateIndex
CREATE UNIQUE INDEX "_PipelineStepDependencies_AB_unique" ON "_PipelineStepDependencies"("A", "B");

-- CreateIndex
CREATE INDEX "_PipelineStepDependencies_B_index" ON "_PipelineStepDependencies"("B");

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_producedById_fkey" FOREIGN KEY ("producedById") REFERENCES "ToolRun"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ToolRun" ADD CONSTRAINT "ToolRun_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ToolRun" ADD CONSTRAINT "ToolRun_pipelineRunId_fkey" FOREIGN KEY ("pipelineRunId") REFERENCES "PipelineRun"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ToolRun" ADD CONSTRAINT "ToolRun_pipelineStepId_fkey" FOREIGN KEY ("pipelineStepId") REFERENCES "PipelineStep"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ContentConsumedBy" ADD CONSTRAINT "_ContentConsumedBy_B_fkey" FOREIGN KEY ("B") REFERENCES "ToolRun"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_PipelineStepDependencies" ADD CONSTRAINT "_PipelineStepDependencies_A_fkey" FOREIGN KEY ("A") REFERENCES "PipelineStep"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_PipelineStepDependencies" ADD CONSTRAINT "_PipelineStepDependencies_B_fkey" FOREIGN KEY ("B") REFERENCES "PipelineStep"("id") ON DELETE CASCADE ON UPDATE CASCADE;
