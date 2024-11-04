import { BullModule } from "@nestjs/bullmq";
import { forwardRef, Module } from "@nestjs/common";

import { ContentModule } from "../content/content.module";
import { EmbeddingsModule } from "../embeddings/embeddings.module";
import { LLMModule } from "../llm/llm.module";
import { LoaderModule } from "../loader/loader.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { RunpodModule } from "../runpod/runpod.module";
import { SpeechModule } from "../speech/speech.module";
import { StorageModule } from "../storage/storage.module";
import { TextChunksModule } from "../text-chunks/text-chunks.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { RunProcessor } from "./run.processor";
import { RunRepository } from "./run.repository";
import { RunsController } from "./runs.controller";
import { RunsService } from "./runs.service";

@Module({
  controllers: [RunsController],
  exports: [RunsService],
  imports: [
    PrismaModule,
    WebsocketsModule,
    BullModule.registerQueue({
      defaultJobOptions: {
        attempts: 1,
      },
      name: "run",
    }),

    forwardRef(() => RunpodModule),
    StorageModule.forRoot(),
    WebsocketsModule,
    BullModule.registerQueue({
      defaultJobOptions: {
        attempts: 1,
      },
      name: "tool",
    }),
    OrganizationsModule,
    ContentModule,
    LoaderModule,
    LLMModule,
    EmbeddingsModule,
    SpeechModule,
    TextChunksModule,
  ],

  providers: [RunsService, RunRepository, RunProcessor],
})
export class RunsModule {}
