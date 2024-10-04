/*
  Warnings:

  - Added the required column `updatedAt` to the `VectorRecord` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "VectorRecord" ADD COLUMN     "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updatedAt" TIMESTAMP(3) NOT NULL;

-- AddForeignKey
ALTER TABLE "VectorRecord" ADD CONSTRAINT "VectorRecord_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;
