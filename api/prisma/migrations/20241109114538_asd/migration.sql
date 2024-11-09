-- CreateTable
CREATE TABLE "_ContentToLabel" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "_ContentToLabel_AB_unique" ON "_ContentToLabel"("A", "B");

-- CreateIndex
CREATE INDEX "_ContentToLabel_B_index" ON "_ContentToLabel"("B");

-- AddForeignKey
ALTER TABLE "_ContentToLabel" ADD CONSTRAINT "_ContentToLabel_A_fkey" FOREIGN KEY ("A") REFERENCES "Content"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ContentToLabel" ADD CONSTRAINT "_ContentToLabel_B_fkey" FOREIGN KEY ("B") REFERENCES "Label"("id") ON DELETE CASCADE ON UPDATE CASCADE;
