import { IntersectionType, PartialType, PickType } from "@nestjs/swagger";
import { ApiProperty } from "@nestjs/swagger";
import { IsOptional, IsString } from "class-validator";

import { ThreadEntity } from "../entities/thread.entity";

export class CreateThreadDto extends IntersectionType(
  PartialType(PickType(ThreadEntity, ["name"] as const))
) {
  @ApiProperty({
    description:
      "Optional. The id to use as the chat id. If taken, this endpoint will return a 409.",
    example: "chatId1",
    required: false,
  })
  @IsOptional()
  @IsString()
  id?: string;
}
