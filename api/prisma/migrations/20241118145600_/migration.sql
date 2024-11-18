/*
  Warnings:

  - A unique constraint covering the columns `[pipelineRunId,pipelineStepId]` on the table `Run` will be added. If there are existing duplicate values, this will fail.

*/
-- DropForeignKey
ALTER TABLE "Run" DROP CONSTRAINT "Run_pipelineId_fkey";

-- AlterTable
ALTER TABLE "Run" ALTER COLUMN "pipelineId" DROP NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "Run_pipelineRunId_pipelineStepId_key" ON "Run"("pipelineRunId", "pipelineStepId");

-- AddForeignKey
ALTER TABLE "Run" ADD CONSTRAINT "Run_pipelineId_fkey" FOREIGN KEY ("pipelineId") REFERENCES "Pipeline"("id") ON DELETE SET NULL ON UPDATE CASCADE;
