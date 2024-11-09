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
import { TransformationProcessor } from "./transformation.processor";
import { TransformationRepository } from "./transformation.repository";
import { TransformationsService } from "./transformations.service";

@Module({
  exports: [TransformationsService],
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
  providers: [
    TransformationsService,
    TransformationRepository,
    TransformationProcessor,
  ],
})
export class TransformationsModule {}
