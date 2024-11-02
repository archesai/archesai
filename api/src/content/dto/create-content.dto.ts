import { ApiProperty, PickType } from "@nestjs/swagger";
import { IsArray } from "class-validator";

import { ContentEntity } from "../entities/content.entity";

export class CreateContentDto extends PickType(ContentEntity, [
  "name",
  "url",
] as const) {
  @ApiProperty({
    description: "The tool IDs to run with this content",
    example: ["tool-id-1", "tool-id-2"],
  })
  @IsArray({
    always: true,
  })
  toolIds: string[];
}
