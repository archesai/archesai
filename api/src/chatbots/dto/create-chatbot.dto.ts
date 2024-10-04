import { PickType } from "@nestjs/swagger";

import { ChatbotEntity } from "../entities/chatbot.entity";

export class CreateChatbotDto extends PickType(ChatbotEntity, [
  "name",
  "description",
  "llmBase",
] as const) {}
