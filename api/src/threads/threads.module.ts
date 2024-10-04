import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { ThreadRepository } from "./thread.repository";
import { ThreadsController } from "./threads.controller";
import { ThreadsCron } from "./threads.cron";
import { ThreadsService } from "./threads.service";

@Module({
  controllers: [ThreadsController],
  exports: [ThreadsService],
  imports: [WebsocketsModule, PrismaModule],
  providers: [ThreadsService, ThreadRepository, ThreadsCron],
})
export class ThreadsModule {}
