import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { JobRepository } from "./job.repository";
import { JobsController } from "./jobs.controller";
import { JobsService } from "./jobs.service";

@Module({
  controllers: [JobsController],
  exports: [JobsService],
  imports: [PrismaModule, WebsocketsModule],
  providers: [JobsService, JobRepository],
})
export class JobsModule {}
