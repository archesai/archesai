import { PickType } from "@nestjs/swagger";

import { MessageEntity } from "../entities/message.entity";

export class CreateMessageDto extends PickType(MessageEntity, [
  "question",
  "answerLength",
  "contextLength",
  "topK",
  "similarityCutoff",
  "temperature",
] as const) {
  answerLength: number = 240;
  contextLength: number = 1000;
  similarityCutoff: number = 0.7;
  temperature: number = 0.7;
  topK: number = 5;
}
