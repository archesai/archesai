import { IntersectionType, PartialType, PickType } from "@nestjs/swagger";
import { ApiProperty } from "@nestjs/swagger";
import { IsOptional, IsString } from "class-validator";

import { ThreadEntity } from "../entities/thread.entity";

export class CreateThreadDto extends IntersectionType(
  PartialType(PickType(ThreadEntity, ["name", "chatbotId"] as const))
) {
  @ApiProperty({
    description:
      "Optional. The id to use as the thread id. If taken, this endpoint will return a 409.",
    example: "custom-thread-id",
    required: false,
  })
  @IsOptional()
  @IsString()
  id?: string;
}
