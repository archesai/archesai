import { PickType } from "@nestjs/swagger";

import { MessageEntity } from "../entities/message.entity";

export class CreateMessageDto extends PickType(MessageEntity, [
  "question",
] as const) {}
