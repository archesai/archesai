import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { PipelineRepository } from "./pipeline.repository";
import { PipelinesController } from "./pipelines.controller";
import { PipelinesService } from "./pipelines.service";

@Module({
  controllers: [PipelinesController],
  exports: [PipelinesService],
  imports: [PrismaModule, WebsocketsModule],
  providers: [PipelinesService, PipelineRepository],
})
export class PipelinesModule {}
