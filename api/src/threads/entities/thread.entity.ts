import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Message, Thread } from "@prisma/client";
import { Organization } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsOptional } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";

export class ThreadEntity extends BaseEntity implements Thread {
  @ApiHideProperty()
  @Exclude()
  chatbotId: string;

  // Public Properties
  @ApiProperty({
    description: "The total number of credits used in this chat",
    example: 10000,
  })
  @Expose()
  credits: number;

  @ApiHideProperty()
  @Exclude()
  messages: Message[];

  @ApiProperty({
    default: "New Chat",
    description: "The chat thread name",
    example: "What are the morals of the story in Aesop's Fables?",
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

  // Private Properties
  @ApiHideProperty()
  @Exclude()
  organization: Organization;

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
