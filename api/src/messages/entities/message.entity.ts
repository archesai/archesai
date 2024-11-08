import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Message } from "@prisma/client";
import { Expose } from "class-transformer";
import { IsString } from "class-validator";

import { BaseEntity } from "../../common/dto/base.entity.dto";

export class MessageEntity extends BaseEntity implements Message {
  @ApiProperty({
    description: "The answer given by the bot",
    example: "The name of this document is Aesop's Fables",
  })
  @Expose()
  answer: string;

  @ApiHideProperty()
  @Expose()
  orgname: string;

  @ApiProperty({
    description: "The question in this message",
    example: "What is the name of this document?",
  })
  @Expose()
  @IsString()
  question: string;

  @ApiProperty({
    description: "The id of the thread this message belongs to",
    example: "thread1",
  })
  @Expose()
  threadId: string;

  constructor(message: Message) {
    super();
    Object.assign(this, message);
  }
}
