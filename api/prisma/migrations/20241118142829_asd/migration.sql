/*
  Warnings:

  - Added the required column `runType` to the `Run` table without a default value. This is not possible if the table is not empty.

*/
-- CreateEnum
CREATE TYPE "RunType" AS ENUM ('PIPELINE_RUN', 'TOOL_RUN');

-- AlterTable
ALTER TABLE "Run" ADD COLUMN     "runType" "RunType" NOT NULL;
