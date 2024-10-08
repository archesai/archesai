import { Module } from "@nestjs/common";

import { ChatbotsModule } from "../chatbots/chatbots.module";
import { ContentModule } from "../content/content.module";
import { EmbeddingsModule } from "../embeddings/embeddings.module";
import { LLMModule } from "../llm/llm.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { ThreadsModule } from "../threads/threads.module";
import { VectorDBModule } from "../vector-db/vector-db.module";
import { VectorRecordModule } from "../vector-records/vector-record.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { MessageRepository } from "./message.repository";
import { MessagesController } from "./messages.controller";
import { MessagesService } from "./messages.service";

@Module({
  controllers: [MessagesController],
  imports: [
    PrismaModule,
    VectorDBModule,
    EmbeddingsModule,
    LLMModule,
    WebsocketsModule,
    OrganizationsModule,
    ChatbotsModule,
    ThreadsModule,
    ContentModule,
    VectorRecordModule,
  ],
  providers: [MessagesService, MessageRepository],
})
export class MessagesModule {}
