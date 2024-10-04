import { PartialType } from "@nestjs/swagger";

import { CreateChatbotDto } from "./create-chatbot.dto";

export class UpdateChatbotDto extends PartialType(CreateChatbotDto) {}
