-- CreateEnum
CREATE TYPE "ARTokenType" AS ENUM ('EMAIL_VERIFICATION', 'PASSWORD_RESET', 'EMAIL_CHANGE');

-- CreateTable
CREATE TABLE "ARToken" (
    "id" TEXT NOT NULL,
    "type" "ARTokenType" NOT NULL,
    "token" TEXT NOT NULL,
    "expiresAt" TIMESTAMP(3) NOT NULL,
    "userId" TEXT NOT NULL,
    "newEmail" TEXT,

    CONSTRAINT "ARToken_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "ARToken" ADD CONSTRAINT "ARToken_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
