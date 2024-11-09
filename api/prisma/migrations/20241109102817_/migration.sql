/*
  Warnings:

  - A unique constraint covering the columns `[pipelineRunId,pipelineStepId]` on the table `Transformation` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "Transformation_pipelineRunId_pipelineStepId_key" ON "Transformation"("pipelineRunId", "pipelineStepId");
