import { Module } from "@nestjs/common";

import { PipelinesModule } from "../pipelines/pipelines.module";
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
    PipelinesModule,
  ],
  providers: [ContentService, ContentRepository],
})
export class ContentModule {}
