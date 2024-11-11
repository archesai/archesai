import {
  ApiProperty,
  IntersectionType,
  PartialType,
  PickType,
} from "@nestjs/swagger";
import { IsArray, IsOptional, IsString } from "class-validator";

import { ContentEntity } from "../entities/content.entity";

export class CreateContentDto extends IntersectionType(
  PickType(ContentEntity, ["name"] as const),
  PartialType(PickType(ContentEntity, ["url", "text"] as const))
) {
  @ApiProperty({
    description: "The labels to associate with the content",
    example: ["label-1", "label-2"],
    required: false,
    type: [String],
  })
  @IsArray()
  @IsOptional()
  @IsString({ each: true })
  labels?: string[] = [];
}
