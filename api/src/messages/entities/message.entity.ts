import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Citation, Message, Thread } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsNumber, IsOptional, IsString } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";
import { CitationEntity } from "../../threads/entities/citation.entity";

export class MessageEntity extends BaseEntity implements Message {
  @ApiProperty({
    description: "The answer given by the bot",
    example: "The name of this document is Aesop's Fables",
  })
  @Expose()
  answer: string;

  @ApiProperty({
    default: 240,
    description: "The max length of the answer given by the bot",
    example: 240,
    required: false,
  })
  @IsOptional()
  @IsNumber()
  @Expose()
  answerLength: number;

  @ApiProperty({
    description: "The sources used in this message",
    type: [CitationEntity],
  })
  @Expose()
  citations: Citation[];

  @ApiProperty({
    default: 1000,
    description: "The max length of the context given to the bot",
    example: 3000,
    required: false,
  })
  @Expose()
  @IsOptional()
  @IsNumber()
  contextLength: number;

  @ApiProperty({
    description: "The number of credits used in this message",
    example: 14,
  })
  @Expose()
  credits: number;

  @ApiProperty({
    description: "The question in this message",
    example: "What is the name of this document?",
  })
  @Expose()
  @IsString()
  question: string;

  @ApiProperty({
    default: 0.7,
    description: "The sililarity cutoff used in this message",
    example: 0.7,
    required: false,
  })
  @IsOptional()
  @Expose()
  @IsOptional()
  @IsNumber()
  similarityCutoff: number;

  @ApiProperty({
    default: 0.7,
    description: "The temperature for the LLM",
    example: 0.7,
    required: false,
  })
  @IsOptional()
  @Expose()
  @IsOptional()
  @IsNumber()
  temperature: number;

  // Private Properties
  @ApiHideProperty()
  @Exclude()
  thread: Thread;

  // Public Properties
  @ApiProperty({
    description: "The id of the thread this message belongs to",
    example: "thread1",
  })
  @Expose()
  threadId: string;

  @ApiProperty({
    default: 5,
    description: "The max number of sources returned included in the context",
    example: 10,
    required: false,
  })
  @Expose()
  @IsNumber()
  @IsOptional()
  topK: number;

  constructor(
    messages: {
      citations: ({
        message: Message;
      } & Citation)[];
    } & Message
  ) {
    super();
    Object.assign(this, messages);
    this.citations = messages.citations.map((s) => new CitationEntity(s));
  }
}
