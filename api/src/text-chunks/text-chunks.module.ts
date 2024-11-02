import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { TextChunkRepository } from "./text-chunk.repository";
import { TextChunksController } from "./text-chunks.controller";
import { TextChunksService } from "./text-chunks.service";

@Module({
  controllers: [TextChunksController],
  exports: [TextChunksService],
  imports: [PrismaModule],
  providers: [TextChunksService, TextChunkRepository],
})
export class TextChunksModule {}
