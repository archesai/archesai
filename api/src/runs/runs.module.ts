import { HttpModule } from "@nestjs/axios";
import { BullModule } from "@nestjs/bullmq";
import { forwardRef, Module } from "@nestjs/common";

import { ContentModule } from "../content/content.module";
import { LLMModule } from "../llm/llm.module";
import { PipelinesModule } from "../pipelines/pipelines.module";
import { PrismaModule } from "../prisma/prisma.module";
import { RunpodModule } from "../runpod/runpod.module";
import { SpeechModule } from "../speech/speech.module";
import { StorageModule } from "../storage/storage.module";
import { ToolsModule } from "../tools/tools.module";
import { RunProcessor } from "./run.processor";
import { RunRepository } from "./run.repository";
import { RunsController } from "./runs.controller";
import { RunsService } from "./runs.service";

@Module({
  controllers: [RunsController],
  exports: [RunsService],
  imports: [
    PrismaModule,
    StorageModule.forRoot(),
    ContentModule,
    LLMModule,
    BullModule.registerQueue({
      name: "run",
    }),
    BullModule.registerFlowProducer({
      name: "flow",
    }),
    SpeechModule,
    HttpModule,
    forwardRef(() => RunpodModule),
    PipelinesModule,
    ToolsModule,
  ],
  providers: [RunsService, RunRepository, RunProcessor],
})
export class RunsModule {}
