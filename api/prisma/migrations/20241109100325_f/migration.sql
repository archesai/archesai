/*
  Warnings:

  - You are about to drop the `RunContent` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "RunContent" DROP CONSTRAINT "RunContent_contentId_fkey";

-- DropForeignKey
ALTER TABLE "RunContent" DROP CONSTRAINT "RunContent_pipelineRunId_fkey";

-- DropForeignKey
ALTER TABLE "RunContent" DROP CONSTRAINT "RunContent_transformationId_fkey";

-- AlterTable
ALTER TABLE "Content" ADD COLUMN     "producedById" TEXT;

-- AlterTable
ALTER TABLE "PipelineStep" ALTER COLUMN "updatedAt" DROP DEFAULT;

-- DropTable
DROP TABLE "RunContent";

-- CreateTable
CREATE TABLE "_ContentConsumedBy" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "_ContentConsumedBy_AB_unique" ON "_ContentConsumedBy"("A", "B");

-- CreateIndex
CREATE INDEX "_ContentConsumedBy_B_index" ON "_ContentConsumedBy"("B");

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_producedById_fkey" FOREIGN KEY ("producedById") REFERENCES "Transformation"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ContentConsumedBy" ADD CONSTRAINT "_ContentConsumedBy_A_fkey" FOREIGN KEY ("A") REFERENCES "Content"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ContentConsumedBy" ADD CONSTRAINT "_ContentConsumedBy_B_fkey" FOREIGN KEY ("B") REFERENCES "Transformation"("id") ON DELETE CASCADE ON UPDATE CASCADE;
