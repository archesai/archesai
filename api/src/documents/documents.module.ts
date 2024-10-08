import { BullModule } from "@nestjs/bull";
import { Module } from "@nestjs/common";

import { ContentModule } from "../content/content.module";
import { EmbeddingsModule } from "../embeddings/embeddings.module";
import { JobsModule } from "../jobs/jobs.module";
import { LLMModule } from "../llm/llm.module";
import { LoaderModule } from "../loader/loader.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { StorageModule } from "../storage/storage.module";
import { VectorRecordModule } from "../vector-records/vector-record.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { DocumentProcessor } from "./document.processor";
import { DocumentsController } from "./documents.controller";
import { DocumentsService } from "./documents.service";

@Module({
  controllers: [DocumentsController],
  imports: [
    BullModule.registerQueue({
      defaultJobOptions: {
        attempts: 1,
      },
      name: "document",
    }),
    ContentModule,
    WebsocketsModule,
    OrganizationsModule,
    JobsModule,
    LoaderModule,
    StorageModule.forRoot(),
    LLMModule,
    EmbeddingsModule,
    VectorRecordModule,
  ],
  providers: [DocumentProcessor, DocumentsService],
})
export class DocumentsModule {}
