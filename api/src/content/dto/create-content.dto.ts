import { ApiProperty, PickType } from "@nestjs/swagger";
import { IsString } from "class-validator";

import { ContentEntity } from "../entities/content.entity";

export class CreateContentDto extends PickType(ContentEntity, [
  "name",
  "url",
] as const) {
  @ApiProperty({
    description: "The tool IDs to run with this content",
    example: ["tool-id-uuid"],
  })
  @IsString()
  pipelineId: string;
}
