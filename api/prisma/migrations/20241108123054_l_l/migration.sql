/*
  Warnings:

  - Added the required column `orgname` to the `Message` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Message" ADD COLUMN     "orgname" TEXT NOT NULL;

-- AddForeignKey
ALTER TABLE "Message" ADD CONSTRAINT "Message_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;
