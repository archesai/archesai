-- DropForeignKey
ALTER TABLE "ApiToken" DROP CONSTRAINT "ApiToken_username_fkey";

-- DropForeignKey
ALTER TABLE "Content" DROP CONSTRAINT "Content_orgname_fkey";

-- DropForeignKey
ALTER TABLE "Member" DROP CONSTRAINT "Member_username_fkey";

-- DropForeignKey
ALTER TABLE "VectorRecord" DROP CONSTRAINT "VectorRecord_contentId_fkey";

-- AddForeignKey
ALTER TABLE "Content" ADD CONSTRAINT "Content_orgname_fkey" FOREIGN KEY ("orgname") REFERENCES "Organization"("orgname") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "VectorRecord" ADD CONSTRAINT "VectorRecord_contentId_fkey" FOREIGN KEY ("contentId") REFERENCES "Content"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Member" ADD CONSTRAINT "Member_username_fkey" FOREIGN KEY ("username") REFERENCES "User"("username") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ApiToken" ADD CONSTRAINT "ApiToken_username_fkey" FOREIGN KEY ("username") REFERENCES "User"("username") ON DELETE CASCADE ON UPDATE CASCADE;
