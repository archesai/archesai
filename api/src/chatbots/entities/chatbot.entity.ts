import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Chatbot } from "@prisma/client";
import { Organization } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsOptional, IsString, ValidateNested } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";
import { ContentFieldItem } from "../../content/dto/content-field-item.dto";

export class ChatbotEntity extends BaseEntity implements Chatbot {
  @ApiProperty({
    description: "The chatbot description",
    example: "You are a chatbot designed to answer questions about Arches AI",
  })
  @Expose()
  @IsString()
  description: string;

  @ApiProperty({
    description: "The documents this chatbot has access to",
    example: [
      {
        id: "417221d3-2b96-4157-a79b-f1866df7743b",
        name: "Arches API Documentation",
      },
    ],
    type: [ContentFieldItem],
  })
  @Expose()
  @IsOptional()
  @ValidateNested({ each: true })
  documents: ContentFieldItem[];

  @ApiProperty({
    default: "GPT_3_5_TURBO_16_K",
    description: "The base LLM that the chatbot will use",
    required: false,
  })
  @Expose()
  @IsOptional()
  llmBase: string;

  @ApiProperty({
    description: "The chatbot name",
    example: "Arches AI Documentation Chatbot",
  })
  @Expose()
  @IsString()
  name: string;

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

  constructor(chatbot: Chatbot) {
    super();
    Object.assign(this, chatbot);
  }
}
