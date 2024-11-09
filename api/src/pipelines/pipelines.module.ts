import { BullModule } from "@nestjs/bullmq";
import { Module } from "@nestjs/common";

import { ContentModule } from "../content/content.module";
import { PrismaModule } from "../prisma/prisma.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { PipelineRepository } from "./pipeline.repository";
import { PipelinesController } from "./pipelines.controller";
import { PipelinesService } from "./pipelines.service";

@Module({
  controllers: [PipelinesController],
  exports: [PipelinesService],
  imports: [
    PrismaModule,
    WebsocketsModule,
    BullModule.registerFlowProducer({
      name: "flow",
    }),
    ContentModule,
  ],
  providers: [PipelinesService, PipelineRepository],
})
export class PipelinesModule {}
