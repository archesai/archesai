import { BullModule } from "@nestjs/bullmq";
import { Module } from "@nestjs/common";

import { ContentModule } from "../content/content.module";
import { EmbeddingsModule } from "../embeddings/embeddings.module";
import { JobsModule } from "../jobs/jobs.module";
import { LLMModule } from "../llm/llm.module";
import { LoaderModule } from "../loader/loader.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { RunpodModule } from "../runpod/runpod.module";
import { SpeechModule } from "../speech/speech.module";
import { StorageModule } from "../storage/storage.module";
import { TextChunksModule } from "../text-chunks/text-chunks.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { ToolProcessor } from "./tool.processor";
import { ToolRepository } from "./tool.repository";
import { ToolsController } from "./tools.controller";
import { ToolsService } from "./tools.service";

@Module({
  controllers: [ToolsController],
  exports: [ToolsService],
  imports: [
    PrismaModule,
    RunpodModule,
    StorageModule.forRoot(),
    WebsocketsModule,
    BullModule.registerQueue({
      defaultJobOptions: {
        attempts: 1,
      },
      name: "tool",
    }),
    OrganizationsModule,
    JobsModule,
    ContentModule,
    LoaderModule,
    LLMModule,
    EmbeddingsModule,
    SpeechModule,
    TextChunksModule,
  ],
  providers: [ToolsService, ToolRepository, ToolProcessor],
})
export class ToolsModule {}
