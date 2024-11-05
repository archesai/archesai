/*
  Warnings:

  - You are about to drop the column `answerLength` on the `Message` table. All the data in the column will be lost.
  - You are about to drop the column `contextLength` on the `Message` table. All the data in the column will be lost.
  - You are about to drop the column `credits` on the `Message` table. All the data in the column will be lost.
  - You are about to drop the column `similarityCutoff` on the `Message` table. All the data in the column will be lost.
  - You are about to drop the column `temperature` on the `Message` table. All the data in the column will be lost.
  - You are about to drop the column `topK` on the `Message` table. All the data in the column will be lost.
  - You are about to drop the `Citation` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `TextChunk` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `_ApiTokenToChatbot` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "Citation" DROP CONSTRAINT "Citation_contentId_fkey";

-- DropForeignKey
ALTER TABLE "Citation" DROP CONSTRAINT "Citation_messageId_fkey";

-- DropForeignKey
ALTER TABLE "TextChunk" DROP CONSTRAINT "TextChunk_contentId_fkey";

-- DropForeignKey
ALTER TABLE "TextChunk" DROP CONSTRAINT "TextChunk_orgname_fkey";

-- DropForeignKey
ALTER TABLE "_ApiTokenToChatbot" DROP CONSTRAINT "_ApiTokenToChatbot_A_fkey";

-- DropForeignKey
ALTER TABLE "_ApiTokenToChatbot" DROP CONSTRAINT "_ApiTokenToChatbot_B_fkey";

-- AlterTable
ALTER TABLE "Content" ADD COLUMN     "embedding" vector,
ADD COLUMN     "text" TEXT;

-- AlterTable
ALTER TABLE "Message" DROP COLUMN "answerLength",
DROP COLUMN "contextLength",
DROP COLUMN "credits",
DROP COLUMN "similarityCutoff",
DROP COLUMN "temperature",
DROP COLUMN "topK";

-- DropTable
DROP TABLE "Citation";

-- DropTable
DROP TABLE "TextChunk";

-- DropTable
DROP TABLE "_ApiTokenToChatbot";
