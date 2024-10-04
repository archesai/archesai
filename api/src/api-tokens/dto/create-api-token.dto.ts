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
      "The ids of the agent this token will have access to. This can not be changed later.",
    example: ["agent1", "agent2"],
    required: false,
  })
  @IsOptional()
  @IsString({ each: true })
  chatbotIds?: string[] = [];
}
