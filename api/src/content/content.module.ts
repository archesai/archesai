import { HttpModule } from "@nestjs/axios";
import { BullModule } from "@nestjs/bull";
import { Module } from "@nestjs/common";

import { AudioModule } from "../audio/audio.module";
import { EmbeddingsModule } from "../embeddings/embeddings.module";
import { JobsModule } from "../jobs/jobs.module";
import { LLMModule } from "../llm/llm.module";
import { LoaderModule } from "../loader/loader.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { RunpodModule } from "../runpod/runpod.module";
import { StorageModule } from "../storage/storage.module";
import { VectorDBModule } from "../vector-db/vector-db.module";
import { ContentController } from "./content.controller";
import { ContentProcessor } from "./content.processor";
import { ContentRepository } from "./content.repository";
import { ContentService } from "./content.service";

@Module({
  controllers: [ContentController],
  exports: [ContentService],
  imports: [
    PrismaModule,
    BullModule.registerQueue({
      defaultJobOptions: {
        attempts: 1,
      },
      name: "content",
    }),
    OrganizationsModule,
    StorageModule.forRoot(),
    HttpModule,
    StorageModule,
    AudioModule,
    JobsModule,
    LoaderModule,
    RunpodModule,
    LLMModule,
    EmbeddingsModule,
    VectorDBModule,
  ],
  providers: [ContentService, ContentRepository, ContentProcessor],
})
export class ContentModule {}
