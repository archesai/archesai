/*
  Warnings:

  - You are about to drop the column `chatbotId` on the `Thread` table. All the data in the column will be lost.
  - You are about to drop the `Chatbot` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "Chatbot" DROP CONSTRAINT "Chatbot_orgname_fkey";

-- DropForeignKey
ALTER TABLE "Thread" DROP CONSTRAINT "Thread_chatbotId_fkey";

-- AlterTable
ALTER TABLE "Thread" DROP COLUMN "chatbotId";

-- DropTable
DROP TABLE "Chatbot";
