import { Module } from "@nestjs/common";

import { ContentModule } from "../content/content.module";
import { PrismaModule } from "../prisma/prisma.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { ThreadRepository } from "./thread.repository";
import { ThreadsController } from "./threads.controller";
import { ThreadsService } from "./threads.service";

@Module({
  controllers: [ThreadsController],
  exports: [ThreadsService],
  imports: [WebsocketsModule, PrismaModule, ContentModule, WebsocketsModule],
  providers: [ThreadsService, ThreadRepository],
})
export class ThreadsModule {}
