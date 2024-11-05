import { Module } from "@nestjs/common";

import { ContentModule } from "../content/content.module";
import { EmbeddingsModule } from "../embeddings/embeddings.module";
import { LLMModule } from "../llm/llm.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { ThreadsModule } from "../threads/threads.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { MessageRepository } from "./message.repository";
import { MessagesController } from "./messages.controller";
import { MessagesService } from "./messages.service";

@Module({
  controllers: [MessagesController],
  imports: [
    PrismaModule,
    EmbeddingsModule,
    LLMModule,
    WebsocketsModule,
    OrganizationsModule,
    ThreadsModule,
    ContentModule,
  ],
  providers: [MessagesService, MessageRepository],
})
export class MessagesModule {}
