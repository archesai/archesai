import { IntersectionType, PartialType, PickType } from "@nestjs/swagger";
import { ApiProperty } from "@nestjs/swagger";
import { IsOptional, IsString } from "class-validator";

import { LabelEntity } from "../entities/label.entity";

export class CreateLabelDto extends IntersectionType(
  PartialType(PickType(LabelEntity, ["name"] as const))
) {
  @ApiProperty({
    description:
      "Optional. The id to use as the label id. If taken, this endpoint will return a 409.",
    example: "custom-label-id",
    required: false,
  })
  @IsOptional()
  @IsString()
  id?: string;
}
