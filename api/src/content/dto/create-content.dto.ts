import { IntersectionType, PartialType, PickType } from "@nestjs/swagger";

import { ContentEntity } from "../entities/content.entity";

export class CreateContentDto extends IntersectionType(
  PickType(ContentEntity, ["name"] as const),
  PartialType(PickType(ContentEntity, ["url", "text"] as const))
) {}
