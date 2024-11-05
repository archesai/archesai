import { ApiProperty } from "@nestjs/swagger";
import { Thread } from "@prisma/client";
import { Expose } from "class-transformer";
import { IsOptional, IsString } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";

export class ThreadEntity extends BaseEntity implements Thread {
  @ApiProperty({
    description: "The chatbot ID",
    example: "fa9023b1-7b7b-4b7b-8b7b-7b7b7b7b7b7b",
    required: false,
  })
  @IsString()
  @IsOptional()
  @Expose()
  chatbotId: string;

  @ApiProperty({
    description: "The total number of credits used in this chat",
    example: 10000,
  })
  @Expose()
  credits: number;

  @ApiProperty({
    default: "New Chat",
    description: "The chat thread name",
    example: "What are the morals of the story in Aesop's Fables?",
    required: false,
  })
  @Expose()
  @IsOptional()
  name: string;

  @ApiProperty({
    description: "The total number of messages in this chat",
    example: 10000,
  })
  @Expose()
  numMessages: number;

  @ApiProperty({
    description: "The organization name",
    example: "my-organization",
  })
  @Expose()
  orgname: string;

  constructor(
    thread: {
      _count: {
        messages: number;
      };
    } & Thread
  ) {
    super();
    Object.assign(this, thread);
    this.numMessages = thread._count.messages;
  }
}
