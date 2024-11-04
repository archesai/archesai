import { PickType } from "@nestjs/swagger";
import { ApiProperty } from "@nestjs/swagger";
import { IsOptional, IsString } from "class-validator";

import { ApiTokenEntity } from "../entities/api-token.entity";

export class CreateApiTokenDto extends PickType(ApiTokenEntity, [
  "role",
  "domains",
  "name",
] as const) {
  @ApiProperty({
    default: [],
    description:
      "The ids of the chatbot this token will have access to. This can not be changed later.",
    example: ["chatbot1", "chatbot2"],
    required: false,
  })
  @IsOptional()
  @IsString({ always: false, each: true })
  chatbotIds?: string[] = [];
}
