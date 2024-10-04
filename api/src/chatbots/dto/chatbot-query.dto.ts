import { ApiProperty } from "@nestjs/swagger";
import { IsOptional, IsString } from "class-validator";

import { SearchQueryDto } from "../../common/search-query";

export class ChatbotQueryDto extends SearchQueryDto {
  @IsOptional()
  @IsString()
  @ApiProperty({
    description: "The name to search for",
    example: "Chatbot 1",
  })
  name?: string;
}
