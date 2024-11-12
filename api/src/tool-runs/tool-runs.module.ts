import { HttpModule } from "@nestjs/axios";
import { BullModule } from "@nestjs/bullmq";
import { forwardRef, Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { ContentModule } from "../content/content.module";
import { LLMModule } from "../llm/llm.module";
import { PrismaModule } from "../prisma/prisma.module";
import { RunpodModule } from "../runpod/runpod.module";
import { SpeechModule } from "../speech/speech.module";
import { StorageModule } from "../storage/storage.module";
import { ToolRunProcessor } from "./tool-run.processor";
import { ToolRunRepository } from "./tool-run.repository";
import { ToolRunsService } from "./tool-runs.service";

@Module({
  exports: [ToolRunsService],
  imports: [
    PrismaModule,
    StorageModule.forRoot(),
    ContentModule,
    LLMModule,
    BullModule.registerQueue({
      name: "flow",
    }),
    SpeechModule,
    HttpModule,
    ConfigModule,
    forwardRef(() => RunpodModule),
  ],
  providers: [ToolRunsService, ToolRunRepository, ToolRunProcessor],
})
export class ToolRunsModule {}
