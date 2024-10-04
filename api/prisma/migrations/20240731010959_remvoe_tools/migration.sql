/*
  Warnings:

  - You are about to drop the `Tool` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `ToolUsage` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `_AgentToTool` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "Tool" DROP CONSTRAINT "Tool_orgname_fkey";

-- DropForeignKey
ALTER TABLE "ToolUsage" DROP CONSTRAINT "ToolUsage_interactionId_fkey";

-- DropForeignKey
ALTER TABLE "ToolUsage" DROP CONSTRAINT "ToolUsage_toolId_fkey";

-- DropForeignKey
ALTER TABLE "_AgentToTool" DROP CONSTRAINT "_AgentToTool_A_fkey";

-- DropForeignKey
ALTER TABLE "_AgentToTool" DROP CONSTRAINT "_AgentToTool_B_fkey";

-- DropTable
DROP TABLE "Tool";

-- DropTable
DROP TABLE "ToolUsage";

-- DropTable
DROP TABLE "_AgentToTool";

-- DropEnum
DROP TYPE "ToolDataType";
