import { BullModule } from "@nestjs/bullmq";
import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { StorageModule } from "../storage/storage.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { ContentController } from "./content.controller";
import { ContentRepository } from "./content.repository";
import { ContentService } from "./content.service";

@Module({
  controllers: [ContentController],
  exports: [ContentService],
  imports: [
    PrismaModule,
    StorageModule.forRoot(),
    WebsocketsModule,
    BullModule.registerQueue({
      defaultJobOptions: {
        attempts: 1,
      },
      name: "tool",
    }),
  ],
  providers: [ContentService, ContentRepository],
})
export class ContentModule {}
