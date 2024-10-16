import { ApiProperty } from "@nestjs/swagger";
import { IsEnum, IsOptional, IsString } from "class-validator";

import { SearchQueryDto } from "../../common/search-query";

export enum SortByField {
  CREATED = "createdAt",
  LLM_BASE = "llmBase",
  NAME = "name",
}

export class ChatbotQueryDto extends SearchQueryDto {
  @IsOptional()
  @IsString()
  @ApiProperty({
    description: "The name to search for",
    example: "Chatbot 1",
  })
  name?: string;

  @ApiProperty({
    default: SortByField.CREATED,
    description: "The field to sort the results by",
    enum: SortByField,
    required: false,
  })
  @IsEnum(SortByField)
  @IsOptional()
  sortBy? = "createdAt";
}
