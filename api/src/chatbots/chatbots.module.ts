import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { ChatbotRepository } from "./chatbot.repository";
import { ChatbotsController } from "./chatbots.controller";
import { ChatbotsService } from "./chatbots.service";

@Module({
  controllers: [ChatbotsController],
  exports: [ChatbotsService],
  imports: [PrismaModule, WebsocketsModule],
  providers: [ChatbotsService, ChatbotRepository],
})
export class ChatbotsModule {}
