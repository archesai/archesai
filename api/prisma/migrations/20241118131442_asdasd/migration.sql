-- AlterTable
ALTER TABLE "ToolRun" ALTER COLUMN "pipelineRunId" DROP NOT NULL,
ALTER COLUMN "pipelineStepId" DROP NOT NULL;
